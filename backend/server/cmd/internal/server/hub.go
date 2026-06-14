package server

import (
	"context"
	"database/sql"
	_ "embed"
	"log"
	"net/http"
	"sailserver/cmd/internal/server/db"
	"sailserver/cmd/internal/server/objects"
	"sailserver/pkg/packets"

	_ "modernc.org/sqlite"
)

//go:embed db/config/schema.sql
var schemaGenSql string

type DbTx struct {
	Ctx     context.Context
	Queries *db.Queries
}

func NewDbTx(h *Hub) *DbTx {
	return &DbTx{
		Ctx:     context.Background(),
		Queries: db.New(h.dbPool),
	}
}

// State handler to process client messages
type ClientStateHandler interface {
	Name() string

	// Inject client into state handler
	SetClient(client ClientInterfacer)

	OnEnter()
	HandleMessage(senderId uint64, message packets.Msg)

	OnExit()
}

type ClientInterfacer interface {
	Id() uint64
	ProcessMessage(senderId uint64, message packets.Msg)

	// Set the client id
	Initialize(id uint64)

	SetState(newState ClientStateHandler)

	SocketSend(message packets.Msg)

	SocketSendAs(message packets.Msg, senderId uint64)

	PassToPeer(message packets.Msg, senderId uint64)

	Broadcast(message packets.Msg)

	ReadPump()
	WritePump()

	DbTx() *DbTx

	Close(reason string)
}

type Hub struct {
	Clients *objects.SharedCollection[ClientInterfacer]

	// Inbound messages from users' clients
	BroadcastChan chan *packets.Packet

	// Processed by all connected clients
	RegisterChan chan ClientInterfacer

	// Unregister channel
	UnregisterChan chan ClientInterfacer

	dbPool *sql.DB
}

func (h *Hub) NewDbTx() {
	panic("unimplemented")
}

func NewHub() *Hub {
	dbPool, err := sql.Open("sqlite", "db.sqlite")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	return &Hub{
		Clients:        objects.NewSharedCollection[ClientInterfacer](),
		BroadcastChan:  make(chan *packets.Packet),
		RegisterChan:   make(chan ClientInterfacer),
		UnregisterChan: make(chan ClientInterfacer),
		dbPool:         dbPool,
	}
}

func (h *Hub) Run() {
	log.Println("Initialising database")
	if _, err := h.dbPool.ExecContext(context.Background(), schemaGenSql); err != nil {
		log.Fatalf("Error initialising database: %v", err)
	}

	log.Println("Awaiting client registrations")
	for {
		select {
		case client := <-h.RegisterChan:
			client.Initialize(h.Clients.Add(client))
		case client := <-h.UnregisterChan:
			h.Clients.Remove(client.Id())
		case packet := <-h.BroadcastChan:
			h.Clients.ForEach(func(clientId uint64, client ClientInterfacer) {
				if clientId != packet.SenderId {
					client.ProcessMessage(packet.SenderId, packet.Msg)
				}
			})
		}
	}
}

func (h *Hub) Serve(getNewClient func(*Hub, http.ResponseWriter, *http.Request) (ClientInterfacer, error), writer http.ResponseWriter, request *http.Request) {
	log.Println("New client connected from ", request.RemoteAddr)
	client, err := getNewClient(h, writer, request)

	if err != nil {
		log.Printf("Error obtaining client for new collection: %v ", err)
		return
	}

	h.RegisterChan <- client

	go client.WritePump()
	go client.ReadPump()
}

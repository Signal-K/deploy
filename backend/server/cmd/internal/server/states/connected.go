package states

import (
	"fmt"
	"log"
	"sailserver/cmd/internal/server"
	"sailserver/pkg/packets"
)

type Connected struct {
	client server.ClientInterfacer
	logger *log.Logger
}

func (c *Connected) Name() string {
	return "Connected"
}

func (c *Connected) SetClient(client server.ClientInterfacer) {
	c.client = client
	loggingPrefix := fmt.Sprintf("Client %d [%s]: ", client.Id(), c.Name())
	c.logger = log.New(log.Writer(), loggingPrefix, log.LstdFlags)
}

func (c *Connected) OnEnter() {
	c.client.SocketSend(packets.NewId(c.client.Id()))

	// Check if user already exists
	// existingUser, err := c.client.DbTx().Queries.GetUserByUsername(c.client.DbTx().Ctx, "test")
	// if err != nil {
	// 	// User doesn't exist, create new user
	// 	user, createErr := c.client.DbTx().Queries.CreateUser(c.client.DbTx().Ctx, db.CreateUserParams{
	// 		Username:     "test",
	// 		PasswordHash: "test1",
	// 	})

	// 	if createErr != nil {
	// 		c.logger.Printf("Failed to create user: %s", createErr)
	// 	} else {
	// 		c.logger.Printf("Created new user: %v", user)
	// 	}
	// } else {
	// 	c.logger.Printf("User already exists: %v", existingUser)
	// }
}

func (c *Connected) HandleMessage(senderId uint64, message packets.Msg) {
	if senderId == c.client.Id() {
		c.client.Broadcast(message)
	} else {
		c.client.SocketSendAs(message, senderId)
	}
}

func (c *Connected) OnExit() {

}

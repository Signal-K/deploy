# Multi-stage build for Go backend
FROM golang:1.24-alpine AS go-builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY backend/server/go.mod backend/server/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY backend/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Final stage - minimal runtime image
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=go-builder /app/main .

# Expose port 8080
EXPOSE 8080

# Run the binary
CMD ["./main"]

// Package transport provides the transport layer for MCP communication.
package transport

import (
	"context"
)

// Message represents a message received from the transport with optional metadata.
type Message struct {
	Data     []byte            // The raw message data
	Metadata map[string]string // Optional metadata (e.g., headers)
}

// RequestHandler is a function that processes a request and returns a response.
// Used by HTTP transport for synchronous request handling.
type RequestHandler func(ctx context.Context, data []byte, metadata map[string]string) []byte

// Transport defines the interface for MCP transport mechanisms.
type Transport interface {
	// Start initializes and starts the transport layer.
	Start(ctx context.Context) error

	// Send sends a message through the transport.
	Send(message []byte) error

	// Receive waits for and returns the next message from the transport.
	// This method blocks until a message is available or an error occurs.
	Receive() ([]byte, error)

	// ReceiveWithMetadata waits for and returns the next message with metadata.
	// Returns the message data and any associated metadata (e.g., headers).
	ReceiveWithMetadata() (*Message, error)

	// Close gracefully shuts down the transport.
	Close() error
}

// SyncTransport extends Transport with synchronous request handling capability.
type SyncTransport interface {
	Transport
	// SetRequestHandler sets the handler for synchronous request processing.
	// When set, requests are processed directly without going through the message queue.
	SetRequestHandler(handler RequestHandler)
}

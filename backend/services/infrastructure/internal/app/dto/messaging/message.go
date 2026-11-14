package messaging

import "context"

// Message represents an event received from a transport.
type Message struct {
    ID string
    Payload []byte
    Metadata map[string]string
}

type MessageHandler func(ctx context.Context, message Message) error

// Event represents an event waiting for durable delivery to the broker.
type Event struct {
    EventID string
    Topic string
    Payload []byte
    Metadata map[string]string
}

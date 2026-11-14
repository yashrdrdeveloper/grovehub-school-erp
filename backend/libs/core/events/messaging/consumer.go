package messaging

import "context"

// Message is a transport-neutral delivery for an event handler.
type Message struct {
    ID string
    Payload []byte
    Metadata map[string]string
}

type Handler func(context.Context, Message) error

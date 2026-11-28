package messaging

import "context"

// Message represents a broker delivery, independent of NATS/Watermill types.
type Message struct {
    ID string
    Payload []byte
    Metadata map[string]string
}

type Handler func(context.Context, Message) error

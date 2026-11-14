package messaging

import (
    "context"

    coreevents "grovehub/libs/core/events"
)

// Message is a broker-delivered message, not durable outbox storage.
type Message struct {
    ID string
    Payload []byte
    Metadata map[string]string
}

type MessageHandler func(ctx context.Context, message Message) error

// Event is a transport-neutral, immutable pending event.
type Event = coreevents.Message

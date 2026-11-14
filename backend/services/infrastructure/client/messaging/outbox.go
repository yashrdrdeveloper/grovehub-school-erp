package messaging

import (
    "context"

    coreevents "grovehub/libs/core/events"
)

// Event is a transport-neutral pending event; no private Infrastructure DTO
// is part of the public signature, so other modules can generate typed mocks.
type Event = coreevents.Message

type EventStore interface {
    Enqueue(ctx context.Context, event coreevents.Message) error
}

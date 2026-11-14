package messaging

import (
    "context"

    messagingdto "grovehub/services/infrastructure/internal/app/dto/messaging"
)

// Event is an immutable outbox command; it is not a broker delivery.
type Event = messagingdto.Event

// EventStore participates in the same transaction as the calling service.
// It is a narrow public contract; claim/retry/cleanup are Infrastructure-only.
type EventStore interface {
    Enqueue(ctx context.Context, event Event) error
}

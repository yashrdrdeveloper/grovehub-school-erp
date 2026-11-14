// Package producers declares Person-owned outbound event contracts.
package person

import (
    "context"

    foundationevents "grovehub/services/foundation/internal/app/constants/events"
    persondto "grovehub/services/foundation/client/person/dto"
)

type ChangedProducer interface {
    Enqueue(ctx context.Context, operation foundationevents.Operation, old, next *persondto.Person) error
}

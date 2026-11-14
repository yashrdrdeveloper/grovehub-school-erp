// Package producers declares Person-owned outbound event contracts.
package producers

import (
    "context"

    foundationevents "grovehub/services/foundation/constants/events"
    persondto "grovehub/services/foundation/dto/person"
)

type ChangedProducer interface {
    Enqueue(ctx context.Context, operation foundationevents.Operation, old, next *persondto.Person) error
}

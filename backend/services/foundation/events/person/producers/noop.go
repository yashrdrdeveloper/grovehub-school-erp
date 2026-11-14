package producers

import (
    "context"

    foundationevents "grovehub/services/foundation/constants/events"
    persondto "grovehub/services/foundation/dto/person"
)

// Noop preserves isolated service tests without messaging. Production injects the outbox producer.
func Noop() ChangedProducer { return noopProducer{} }

type noopProducer struct{}

func (noopProducer) Enqueue(context.Context, foundationevents.Operation, *persondto.Person, *persondto.Person) error { return nil }

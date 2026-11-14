package helpers

import (
    "context"
    "time"

    coreevents "grovehub/libs/core/events"
    foundationevents "grovehub/services/foundation/constants/events"
    persondto "grovehub/services/foundation/dto/person"
    personevents "grovehub/services/foundation/events/person"
)

// BuildChanged uses the stable versioned Person Changed event envelope.
func BuildChanged(ctx context.Context, eventID string, operation foundationevents.Operation, old, next *persondto.Person) coreevents.Envelope[personevents.Changed] {
    gid := ""
    if next != nil { gid = next.GID } else if old != nil { gid = old.GID }
    return coreevents.NewEnvelope(ctx, coreevents.Metadata{
        EventID: eventID,
        EventType: personevents.ChangedType,
        Version: personevents.ChangedVersion,
        AggregateID: gid,
        OccurredAt: time.Now().UTC(),
    }, personevents.Changed{Operation: operation, Old: old, New: next})
}

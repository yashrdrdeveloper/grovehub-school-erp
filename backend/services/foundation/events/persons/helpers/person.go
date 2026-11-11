package helpers

import (
    "context"
    "time"

    coreevents "grovehub/libs/core/events"
    foundationevents "grovehub/services/foundation/constants/events"
    foundationresponse "grovehub/services/foundation/dto/response/person"
    personevents "grovehub/services/foundation/events/persons"
)

// BuildChanged uses the shared generic envelope; the domain owns only its data.
func BuildChanged(
    ctx context.Context,
    eventID string,
    operation foundationevents.Operation,
    old, next *foundationresponse.Person,
) coreevents.Envelope[personevents.Changed] {
    gid := ""
    if next != nil {
        gid = next.GID
    } else if old != nil {
        gid = old.GID
    }
    return coreevents.NewEnvelope(ctx, coreevents.Metadata{
        EventID:     eventID,
        EventType:   personevents.ChangedType,
        Version:     personevents.ChangedVersion,
        AggregateID: gid,
        OccurredAt:  time.Now().UTC(),
    }, personevents.Changed{Operation: operation, Old: old, New: next})
}

func ChangedMetadata(eventID string) map[string]string {
    return map[string]string{
        "event-id":   eventID,
        "event-type": personevents.ChangedType,
        "source":     personevents.Source,
    }
}

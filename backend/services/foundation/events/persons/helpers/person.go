package helpers

import (
    "time"

    coreevents "grovehub/libs/core/events"
    foundationevents "grovehub/services/foundation/constants/events"
    foundationresponse "grovehub/services/foundation/dto/response"
    personevents "grovehub/services/foundation/events/persons"
)

func BuildChanged(eventID string, operation foundationevents.Operation, old, next *foundationresponse.Person) coreevents.Envelope[personevents.Changed] {
    return coreevents.Envelope[personevents.Changed]{
        EventID: eventID,
        EventType: personevents.ChangedType,
        Source: personevents.Source,
        OccurredAt: time.Now().UTC(),
        Data: personevents.Changed{Operation: operation, Old: old, New: next},
    }
}

func ChangedMetadata(eventID string) map[string]string {
    return map[string]string{
        "event-id": eventID,
        "event-type": personevents.ChangedType,
        "source": personevents.Source,
    }
}

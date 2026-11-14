package helpers

import (
    personevents "grovehub/services/foundation/internal/app/events/producers/person/operations/changed"
)

func ChangedMetadata(eventID string) map[string]string {
    return map[string]string{
        "event-id": eventID,
        "event-type": personevents.ChangedType,
        "source": personevents.Source,
    }
}

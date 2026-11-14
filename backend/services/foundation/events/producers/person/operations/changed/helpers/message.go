package helpers

import (
    "context"
    "encoding/json"

    "github.com/google/uuid"

    coreoutbox "grovehub/libs/core/events/outbox"
    personevents "grovehub/services/foundation/events/producers/person/operations/changed"
    eventflows "grovehub/services/foundation/events/producers/person/operations/changed/flows"
)

// BuildMessage creates the immutable bytes and metadata stored by the outbox.
func BuildMessage(ctx context.Context, change eventflows.Change) (coreoutbox.Message, error) {
    eventID := uuid.NewString()
    envelope := BuildChanged(ctx, eventID, change.Operation, change.Old, change.New)
    payload, err := json.Marshal(envelope)
    if err != nil { return coreoutbox.Message{}, err }
    return coreoutbox.Message{
        EventID: eventID,
        Topic: personevents.ChangedTopic,
        Payload: payload,
        Metadata: ChangedMetadata(eventID),
    }, nil
}

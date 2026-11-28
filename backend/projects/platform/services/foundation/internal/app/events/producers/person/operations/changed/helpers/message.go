package helpers

import (
    "context"
    "encoding/json"

    "github.com/google/uuid"

    messagingclient "grovehub/services/infrastructure/client/messaging"
    personevents "grovehub/services/foundation/internal/app/events/producers/person/operations/changed"
    eventflows "grovehub/services/foundation/internal/app/events/producers/person/operations/changed/flows"
)

// BuildMessage creates the immutable bytes and metadata stored by the outbox.
func BuildMessage(ctx context.Context, change eventflows.Change) (messagingclient.Event, error) {
    eventID := uuid.NewString()
    envelope := BuildChanged(ctx, eventID, change.Operation, change.Old, change.New)
    payload, err := json.Marshal(envelope)
    if err != nil { return messagingclient.Event{}, err }
    return messagingclient.Event{
        EventID: eventID,
        Topic: personevents.ChangedTopic,
        Payload: payload,
        Metadata: ChangedMetadata(eventID),
    }, nil
}

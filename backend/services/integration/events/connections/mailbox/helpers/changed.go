package helpers

import (
 "context"
 "time"

 "github.com/google/uuid"

 coreevents "grovehub/libs/core/events"
 mailboxevents "grovehub/services/integration/events/connections/mailbox"
)

func BuildChanged(ctx context.Context,payload mailboxevents.Payload)mailboxevents.Changed{
 return coreevents.NewEnvelope(ctx,coreevents.Metadata{
  EventID:uuid.NewString(),
  EventType:mailboxevents.ChangedEventType,
  Version:1,
  AggregateID:payload.ConnectionID,
  OccurredAt:time.Now().UTC(),
 },payload)
}

func ChangedMetadata(event mailboxevents.Changed)map[string]string{
 return map[string]string{
  "eventId":event.EventID,
  "eventType":event.EventType,
  "source":"integration",
 }
}

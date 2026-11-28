package impl

import (
	"context"
	"encoding/json"
	"errors"

	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	coreevents "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events"
	client "grovehub/services/audit/client/producer"
	"grovehub/services/audit/internal/app/constants/events"
	"grovehub/services/audit/internal/app/events/consumers/activity/helpers"
	"grovehub/services/audit/internal/app/services/producer"
	"grovehub/services/infrastructure/client/messaging"
)

type Service struct {
	store messaging.EventStore
}

func New(store messaging.EventStore) producer.Service {
	return &Service{store:store}
}

func (service *Service) Produce(ctx context.Context, event client.Event) error {
	if service.store == nil || !transaction.Active(ctx) {
		return errors.New("audit producer requires a transactional outbox")
	}
	topic, known := events.Topic(event.Type)
	if !known { return helpers.ErrInvalidEvent }
	entry, err := helpers.Project(ctx, event.Type, event.Payload)
	if err != nil { return err }
	var envelope coreevents.Envelope[json.RawMessage]
	if err := json.Unmarshal(event.Payload, &envelope); err != nil { return helpers.ErrInvalidEvent }
	return service.store.Enqueue(ctx, messaging.Event{EventID:envelope.EventID, Topic:topic, Payload:event.Payload, Metadata:map[string]string{
		"event-id":envelope.EventID,"event-type":envelope.EventType,"source":entry.Source,"correlation-id":envelope.CorrelationID,"causation-id":envelope.CausationID,
	}})
}

func (service *Service) Record(ctx context.Context, activity client.Activity) error {
	if activity.AssetVersion>0 && activity.OccurrenceID=="" { activity.OccurrenceID=helpers.VersionEventID(activity) }
	envelope:=coreevents.NewEnvelope(ctx,coreevents.Metadata{EventID:activity.OccurrenceID,EventType:client.ActivityOccurred,Version:1,AggregateID:activity.AssetID,OccurredAt:activity.OccurredAt},activity)
	payload,err:=json.Marshal(envelope)
	if err!=nil { return err }
	return service.Produce(ctx,client.Event{Type:client.ActivityOccurred,Payload:payload})
}

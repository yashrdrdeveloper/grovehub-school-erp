package impl

import (
	"context"
	"encoding/json"

	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	coreevents "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events"
	"grovehub/services/identity/internal/app/entity"
	producer "grovehub/services/identity/internal/app/events/producers/registration"
	repo "grovehub/services/identity/internal/app/repository/account"
	messaging "grovehub/services/infrastructure/client/messaging"
)

const EventType = "identity.account.registration.changed.v1"
const Topic = "identity_account_registration_changed"

type Snapshot struct {
	AccountID string `json:"accountId"`
	PersonGID string `json:"personGid"`
	Status entity.AccountStatus `json:"status"`
	Version int64 `json:"version"`
}

type Changed struct {
	TenantID string `json:"tenantId"`
	ActorID string `json:"actorId"`
	Operation string `json:"operation"`
	Old *Snapshot `json:"old"`
	New Snapshot `json:"new"`
}

type Producer struct {
	store messaging.EventStore
}

func New(store messaging.EventStore) producer.Producer {
	return &Producer{store: store}
}

func (producer *Producer) Enqueue(ctx context.Context, registration repo.Registration) error {
	if producer.store == nil || !transaction.Active(ctx) {
		return repo.ErrConflict
	}
	data := Changed{TenantID: registration.TenantID, ActorID: registration.ActorID, Operation: "CREATED", New: Snapshot{AccountID: registration.AccountID, PersonGID: registration.PersonGID, Status: entity.PendingSetup, Version: 1}}
	envelope := coreevents.NewEnvelope(ctx, coreevents.Metadata{EventID: registration.EventID, EventType: EventType, Version: 1, AggregateID: registration.AccountID, OccurredAt: registration.At}, data)
	if envelope.CorrelationID == "" {
		envelope.CorrelationID = envelope.EventID
	}
	payload, err := json.Marshal(envelope)
	if err != nil {
		return err
	}
	return producer.store.Enqueue(ctx, messaging.Event{EventID: envelope.EventID, Topic: Topic, Payload: payload, Metadata: map[string]string{"event-id": envelope.EventID, "event-type": EventType, "source": "identity", "correlation-id": envelope.CorrelationID, "causation-id": envelope.CausationID}})
}

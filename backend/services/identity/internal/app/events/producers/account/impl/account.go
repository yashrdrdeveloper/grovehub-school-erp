package impl

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	coreevents "grovehub/libs/core/events"
	"grovehub/services/identity/internal/app/entity"
	producer "grovehub/services/identity/internal/app/events/producers/account"
	changed "grovehub/services/identity/internal/app/events/producers/account/operations/changed"
	infrastructure "grovehub/services/infrastructure/client/messaging"
)

type Producer struct {
	events infrastructure.EventStore
}

var _ producer.ChangedProducer = (*Producer)(nil)

func New(events infrastructure.EventStore) producer.ChangedProducer {
	return &Producer{events: events}
}

func snapshot(account *entity.Account) *changed.Snapshot {
	if account == nil {
		return nil
	}
	tenantID := ""
	if account.TenantID != nil {
		tenantID = *account.TenantID
	}
	return &changed.Snapshot{
		AccountID: account.AccountID,
		PersonGID: account.PersonGID,
		TenantID: tenantID,
		Username: account.Username,
		Status: string(account.Status),
		Version: account.Version,
	}
}

func (publisher *Producer) Enqueue(ctx context.Context, operation string, old *entity.Account, next *entity.Account) error {
	eventID := uuid.NewString()
	accountID := ""
	if next != nil {
		accountID = next.AccountID
	} else if old != nil {
		accountID = old.AccountID
	}

	envelope := coreevents.NewEnvelope(ctx, coreevents.Metadata{
		EventID: eventID,
		EventType: "identity.account.changed.v1",
		Version: 1,
		AggregateID: accountID,
		OccurredAt: time.Now().UTC(),
	}, changed.Changed{
		Operation: operation,
		Old: snapshot(old),
		New: snapshot(next),
	})
	payload, err := json.Marshal(envelope)
	if err != nil {
		return err
	}
	return publisher.events.Enqueue(ctx, infrastructure.Event{
		EventID: eventID,
		Topic: "identity_account_changed",
		Payload: payload,
		Metadata: map[string]string{
			"event-id": eventID,
			"event-type": "identity.account.changed.v1",
			"source": "identity",
		},
	})
}

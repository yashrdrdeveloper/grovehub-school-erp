package impl

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/google/uuid"

	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	coreevents "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events"
	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/events/producers/accountmetadata"
	"grovehub/services/identity/internal/app/events/producers/accountmetadata/operations/changed"
	"grovehub/services/identity/internal/app/repository/account"
	"grovehub/services/identity/internal/app/repository/accountmetadataaudit"
	"grovehub/services/identity/internal/app/services/security/principal"
	messaging "grovehub/services/infrastructure/client/messaging"
)

const EventType = "identity.account.metadata.changed.v1"
const Topic = "identity_account_metadata_changed"

type Recorder struct {
	audit accountmetadataaudit.Writer
	events messaging.EventStore
}

func New(audit accountmetadataaudit.Writer, events messaging.EventStore) accountmetadata.Recorder {
	return &Recorder{audit: audit, events: events}
}

func (recorder *Recorder) Record(ctx context.Context, before entity.Account, after entity.Account) error {
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return err
	}
	if recorder.audit == nil || recorder.events == nil || !transaction.Active(ctx) || before.TenantID == nil || after.TenantID == nil || *before.TenantID != actor.TenantID || *after.TenantID != actor.TenantID || before.AccountID == "" || before.AccountID != after.AccountID || before.Version < 1 || after.Version <= before.Version || after.Version-before.Version != 1 || after.UpdatedAt.IsZero() || after.Username == nil || reflect.DeepEqual(before.Username, after.Username) {
		return account.ErrConflict
	}
	// Metadata changes must not change the projected identity or authority.
	expected := before
	expected.Username = after.Username
	expected.Version = after.Version
	expected.UpdatedAt = after.UpdatedAt
	expected.SearchDetails = after.SearchDetails
	if !reflect.DeepEqual(expected, after) {
		return account.ErrConflict
	}
	data := changed.Changed{TenantID: actor.TenantID, ActorID: actor.AccountID, Operation: "UPDATED", Old: changed.Snapshot{AccountID: before.AccountID, Username: before.Username, Version: before.Version}, New: changed.Snapshot{AccountID: after.AccountID, Username: after.Username, Version: after.Version}}
	envelope := coreevents.NewEnvelope(ctx, coreevents.Metadata{EventID: uuid.NewString(), EventType: EventType, Version: 1, AggregateID: after.AccountID, OccurredAt: after.UpdatedAt}, data)
	if envelope.CorrelationID == "" {
		envelope.CorrelationID = envelope.EventID
	}
	payload, err := json.Marshal(envelope)
	if err != nil {
		return err
	}
	if err := recorder.audit.Insert(ctx, accountmetadataaudit.Entry{EventID: envelope.EventID, TenantID: actor.TenantID, ActorID: actor.AccountID, AccountID: after.AccountID, Version: after.Version, Payload: payload, OccurredAt: after.UpdatedAt}); err != nil {
		return err
	}
	return recorder.events.Enqueue(ctx, messaging.Event{EventID: envelope.EventID, Topic: Topic, Payload: payload, Metadata: map[string]string{"event-id": envelope.EventID, "event-type": EventType, "source": "identity", "correlation-id": envelope.CorrelationID, "causation-id": envelope.CausationID}})
}

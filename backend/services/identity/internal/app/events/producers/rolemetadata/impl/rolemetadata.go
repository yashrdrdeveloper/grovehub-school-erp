package impl

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"time"

	"github.com/google/uuid"

	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	coreevents "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events"
	"grovehub/services/identity/internal/app/entity"
	resource "grovehub/services/identity/internal/app/events/producers/resource/operations/changed"
	"grovehub/services/identity/internal/app/events/producers/rolemetadata"
	"grovehub/services/identity/internal/app/events/producers/rolemetadata/operations/changed"
	"grovehub/services/identity/internal/app/repository/rolemetadataaudit"
	"grovehub/services/identity/internal/app/services/security/principal"
	infrastructure "grovehub/services/infrastructure/client/messaging"
)

const EventType = "identity.role.metadata.changed.v1"
const Topic = "identity_role_metadata_changed"

var ErrInvalid = errors.New("invalid role metadata audit transition")

type Recorder struct {
	audit rolemetadataaudit.Writer
	events infrastructure.EventStore
}

func New(audit rolemetadataaudit.Writer, events infrastructure.EventStore) rolemetadata.Recorder {
	return &Recorder{audit: audit, events: events}
}

func (recorder *Recorder) Record(ctx context.Context, before *resource.Snapshot, after resource.Snapshot) error {
	if recorder.audit == nil || recorder.events == nil || !transaction.Active(ctx) || after.ID == "" || after.Role == nil || after.Account != nil || after.Role.Status != entity.RoleInactive {
		return ErrInvalid
	}
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return err
	}
	operation := "CREATED"
	if before == nil {
		if after.Version != 1 || len(after.Role.Permissions) != 0 {
			return ErrInvalid
		}
	} else {
		operation = "UPDATED"
		if before.ID != after.ID || before.Version < 1 || after.Version <= before.Version || after.Version-before.Version != 1 || before.Role == nil || before.Account != nil || before.Role.Status != entity.RoleInactive || !slices.Equal(before.Role.Permissions, after.Role.Permissions) {
			return ErrInvalid
		}
	}
	at := time.Now().UTC()
	data := changed.Changed{TenantID: actor.TenantID, ActorID: actor.AccountID, ResourceType: "ROLE", Operation: operation, Old: before, New: after}
	envelope := coreevents.NewEnvelope(ctx, coreevents.Metadata{EventID: uuid.NewString(), EventType: EventType, Version: 1, AggregateID: after.ID, OccurredAt: at}, data)
	if envelope.CorrelationID == "" {
		envelope.CorrelationID = envelope.EventID
	}
	payload, err := json.Marshal(envelope)
	if err != nil {
		return err
	}
	if err := recorder.audit.Insert(ctx, rolemetadataaudit.Entry{EventID: envelope.EventID, TenantID: actor.TenantID, ActorID: actor.AccountID, RoleID: after.ID, Version: after.Version, Operation: operation, Payload: payload, OccurredAt: at}); err != nil {
		return err
	}
	return recorder.events.Enqueue(ctx, infrastructure.Event{EventID: envelope.EventID, Topic: Topic, Payload: payload, Metadata: map[string]string{"event-id": envelope.EventID, "event-type": EventType, "source": "identity", "correlation-id": envelope.CorrelationID, "causation-id": envelope.CausationID}})
}

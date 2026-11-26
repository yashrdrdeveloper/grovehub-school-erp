package impl

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"

	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	coreevents "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events"
	kind "grovehub/services/identity/internal/app/constants/approval"
	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/events/producers/resource"
	"grovehub/services/identity/internal/app/events/producers/resource/operations/changed"
	"grovehub/services/identity/internal/app/repository/resourceaudit"
	"grovehub/services/identity/internal/app/services/security/principal"
	infrastructure "grovehub/services/infrastructure/client/messaging"
)

var ErrInvalid = errors.New("invalid resource audit transition")

type Recorder struct {
	audit resourceaudit.Writer
	events infrastructure.EventStore
}

var _ resource.Recorder = (*Recorder)(nil)

func New(audit resourceaudit.Writer, events infrastructure.EventStore) resource.Recorder {
	return &Recorder{audit: audit, events: events}
}

func (recorder *Recorder) Record(ctx context.Context, request entity.ApprovalRequest, before changed.Snapshot, after changed.Snapshot, at time.Time) error {
	if recorder.audit == nil || recorder.events == nil || !transaction.Active(ctx) {
		return ErrInvalid
	}
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return err
	}
	if actor.TenantID != request.TenantID || actor.AccountID != request.RequesterID || request.Status != kind.Approved || before.ID == "" || before.ID != after.ID || before.Version < 1 || after.Version <= before.Version || after.Version-before.Version != 1 || at.IsZero() {
		return ErrInvalid
	}
	resourceType, eventType, topic, err := eventContract(request.RequestType, before, after)
	if err != nil {
		return err
	}
	eventID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	data := changed.Changed{TenantID: request.TenantID, RequestID: request.RequestID, ActorID: actor.AccountID, ResourceType: resourceType, Operation: request.RequestType, Old: before, New: after}
	envelope := coreevents.NewEnvelope(ctx, coreevents.Metadata{EventID: eventID.String(), EventType: eventType, Version: 1, AggregateID: after.ID, OccurredAt: at, CausationID: request.RequestID}, data)
	if envelope.CorrelationID == "" {
		envelope.CorrelationID = request.RequestID
	}
	payload, err := json.Marshal(envelope)
	if err != nil {
		return err
	}
	if err := recorder.audit.Insert(ctx, resourceaudit.Entry{EventID: envelope.EventID, TenantID: request.TenantID, RequestID: request.RequestID, ActorID: actor.AccountID, ResourceType: resourceType, ResourceID: after.ID, Version: after.Version, Operation: string(request.RequestType), Payload: payload, OccurredAt: at}); err != nil {
		return err
	}
	return recorder.events.Enqueue(ctx, infrastructure.Event{EventID: envelope.EventID, Topic: topic, Payload: payload, Metadata: map[string]string{"event-id": envelope.EventID, "event-type": eventType, "source": "identity", "correlation-id": envelope.CorrelationID, "causation-id": envelope.CausationID}})
}

func eventContract(operation kind.RequestType, before changed.Snapshot, after changed.Snapshot) (string, string, string, error) {
	switch operation {
	case kind.RoleLifecycleChange, kind.RolePermissionChange:
		if before.Role != nil && after.Role != nil && before.Account == nil && after.Account == nil {
			return "ROLE", "identity.role.changed.v1", "identity_role_changed", nil
		}
	case kind.RoleAssignmentChange, kind.AccountLifecycleChange:
		if before.Account != nil && after.Account != nil && before.Role == nil && after.Role == nil {
			return "ACCOUNT", "identity.account.authority.changed.v1", "identity_account_authority_changed", nil
		}
	}
	return "", "", "", ErrInvalid
}

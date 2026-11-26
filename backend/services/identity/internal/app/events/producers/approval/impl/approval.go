package impl

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"

	coreevents "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events"
	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	"grovehub/services/identity/internal/app/entity"
	producer "grovehub/services/identity/internal/app/events/producers/approval"
	"grovehub/services/identity/internal/app/events/producers/approval/operations/changed"
	"grovehub/services/identity/internal/app/repository/approvalaudit"
	"grovehub/services/identity/internal/app/services/security/principal"
	infrastructure "grovehub/services/infrastructure/client/messaging"
)

const EventType = "identity.approval.changed.v1"
const Topic = "identity_approval_changed"

var ErrUnavailable = errors.New("approval audit unavailable")

type Recorder struct {
	audit approvalaudit.Writer
	events infrastructure.EventStore
}

var _ producer.Recorder = (*Recorder)(nil)

func New(audit approvalaudit.Writer, events infrastructure.EventStore) producer.Recorder {
	return &Recorder{audit: audit, events: events}
}

func snapshot(request entity.ApprovalRequest) changed.Snapshot {
	return changed.Snapshot{Status: request.Status, Version: request.Version, CurrentStep: request.CurrentStep, StepCount: request.StepCount, ExpiresAt: request.ExpiresAt}
}

func (recorder *Recorder) Record(ctx context.Context, operation string, before *entity.ApprovalRequest, after entity.ApprovalRequest) error {
	if recorder.audit == nil || recorder.events == nil || !transaction.Active(ctx) {
		return ErrUnavailable
	}
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return err
	}
	if actor.TenantID != after.TenantID || after.Version < 1 || after.UpdatedAt.IsZero() {
		return ErrUnavailable
	}
	if before == nil {
		if operation != "CREATED" || after.Version != 1 || actor.AccountID != after.RequesterID {
			return ErrUnavailable
		}
	} else if before.TenantID != after.TenantID || before.RequestID != after.RequestID || after.Version != before.Version+1 {
		return ErrUnavailable
	}
	data := changed.Changed{TenantID: after.TenantID, RequestID: after.RequestID, RequestType: after.RequestType, PolicyID: after.PolicyID, RequesterID: after.RequesterID, ActorID: actor.AccountID, PayloadSHA256: after.PayloadSHA256, Operation: operation, New: snapshot(after)}
	if before != nil {
		old := snapshot(*before)
		data.Old = &old
	}
	eventID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	envelope := coreevents.NewEnvelope(ctx, coreevents.Metadata{EventID: eventID.String(), EventType: EventType, Version: 1, AggregateID: after.RequestID, OccurredAt: after.UpdatedAt, CausationID: after.RequestID}, data)
	if envelope.CorrelationID == "" {
		envelope.CorrelationID = after.RequestID
	}
	payload, err := json.Marshal(envelope)
	if err != nil {
		return err
	}
	if err := recorder.audit.Insert(ctx, approvalaudit.Entry{EventID: envelope.EventID, TenantID: after.TenantID, RequestID: after.RequestID, ActorID: actor.AccountID, Version: after.Version, Operation: operation, Payload: payload, OccurredAt: after.UpdatedAt}); err != nil {
		return err
	}
	return recorder.events.Enqueue(ctx, infrastructure.Event{EventID: envelope.EventID, Topic: Topic, Payload: payload, Metadata: map[string]string{"event-id": envelope.EventID, "event-type": EventType, "source": "identity", "correlation-id": envelope.CorrelationID, "causation-id": envelope.CausationID}})
}

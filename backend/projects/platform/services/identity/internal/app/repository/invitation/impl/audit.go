package impl

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
	coredatabase "grovehub/libs/core/database"
	coreevents "grovehub/libs/core/events"
	auditclient "grovehub/services/audit/client/producer"
	"grovehub/services/identity/internal/app/entity"
)

var deliveryNamespace = uuid.MustParse("60b74248-62f8-4de1-928a-dca56dc32dab")

func (repository *Repository) record(ctx context.Context, entry entity.Invitation, operation string, reason string, next *time.Time, at time.Time) error {
	var scope struct{ TenantID *string }
	if err := coredatabase.TransactionDB(ctx, repository.db).Raw("SELECT tenant_id FROM identity.accounts WHERE account_id = ?", entry.AccountID).Scan(&scope).Error; err != nil {
		return err
	}
	// A legacy orphan or tenantless record has no safe tenant visibility.
	// Do not block suppression of an obsolete job or assign a fabricated tenant.
	if scope.TenantID == nil || *scope.TenantID == "" {
		if operation == "QUEUED" {
			return errors.New("invitation audit requires tenant scope")
		}
		return nil
	}
	if repository.audit == nil {
		return errors.New("invitation audit producer unavailable")
	}
	purpose := entry.Purpose
	if purpose != "ACCOUNT_SETUP" && purpose != "PASSWORD_RESET" {
		purpose = "UNKNOWN"
	}
	eventID := uuid.NewSHA1(deliveryNamespace, []byte(entry.ID+"/"+strconv.Itoa(entry.Attempts)+"/"+operation)).String()
	data := auditclient.InvitationEvent{TenantID: *scope.TenantID, AccountID: entry.AccountID, Purpose: purpose, Operation: operation, Attempt: entry.Attempts, Reason: reason, NextAttemptAt: next}
	event := coreevents.NewEnvelope(ctx, coreevents.Metadata{EventID: eventID, EventType: auditclient.InvitationDelivery, Version: 1, AggregateID: entry.ID, OccurredAt: at}, data)
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return repository.audit.Produce(ctx, auditclient.Event{Type: event.EventType, Payload: payload})
}

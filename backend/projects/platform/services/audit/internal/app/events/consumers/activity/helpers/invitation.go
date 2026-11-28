package helpers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/google/uuid"

	"grovehub/services/audit/internal/app/entity"
	validator "grovehub/services/audit/internal/app/validators/events"
)

func projectInvitation(_ context.Context, raw []byte) (entity.Activity, error) {
	event, err := validator.Invitation(raw)
	if err != nil {
		return entity.Activity{}, err
	}

	data := event.Data
	version := int64(data.Attempt) + 1
	key, _ := json.Marshal([]any{data.TenantID, "INVITATION", event.AggregateID, event.EventType, data.Operation, version})
	id := uuid.NewSHA1(namespace, key).String()

	account := data.AccountID
	tenant := data.TenantID
	data.TenantID = ""
	projected, err := json.Marshal(data)
	if err != nil {
		return entity.Activity{}, ErrInvalidEvent
	}

	semantic, _ := json.Marshal([]any{string(key), json.RawMessage(projected)})
	digest := sha256.Sum256(semantic)

	return entity.Activity{
		AuditID:          id,
		SourceEventID:    event.EventID,
		TenantID:         tenant,
		AssetType:        "INVITATION",
		AssetID:          event.AggregateID,
		AuditType:        event.EventType,
		Operation:        data.Operation,
		Source:           "identity",
		AssetVersion:     version,
		SubjectAccountID: &account,
		Data:             projected,
		SemanticHash:     hex.EncodeToString(digest[:]),
		OccurredAt:       event.OccurredAt.UTC(),
	}, nil
}

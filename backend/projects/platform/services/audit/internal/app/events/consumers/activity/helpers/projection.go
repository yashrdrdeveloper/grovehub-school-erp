package helpers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/google/uuid"

	"grovehub/services/audit/internal/app/constants/events"
	dto "grovehub/services/audit/internal/app/dto/events"
	"grovehub/services/audit/internal/app/entity"
	validator "grovehub/services/audit/internal/app/validators/events"
	validationhelpers "grovehub/services/audit/internal/app/validators/events/helpers"
)

var ErrInvalidEvent = validator.ErrInvalidEvent
var namespace = uuid.MustParse(validationhelpers.NamespaceID)

func optionalID(value string) *string {
	if !validationhelpers.ValidID(value) {
		return nil
	}
	return &value
}

func Project(ctx context.Context, expectedType string, raw []byte) (entity.Activity, error) {
	if expectedType == events.ActivityOccurred {
		return projectOccurrence(raw)
	}
	if expectedType == events.InvitationDelivery {
		return projectInvitation(ctx, raw)
	}

	event, tenant, assetType, err := validator.Identity(expectedType, raw)
	if err != nil {
		return entity.Activity{}, err
	}

	data := event.Data
	for _, value := range []*dto.Snapshot{data.Old, data.New} {
		if value != nil {
			value.TenantID = ""
		}
	}
	projected, err := json.Marshal(struct {
		Old *dto.Snapshot `json:"old"`
		New *dto.Snapshot `json:"new"`
		RequestID *string `json:"requestId,omitempty"`
		RequestType string `json:"requestType,omitempty"`
	}{data.Old, data.New, optionalID(data.RequestID), data.RequestType})
	if err != nil {
		return entity.Activity{}, ErrInvalidEvent
	}
	key, _ := json.Marshal([]any{tenant, assetType, event.AggregateID, expectedType, data.Operation, data.New.Version})
	auditID := uuid.NewSHA1(namespace, key).String()
	semantic, _ := json.Marshal([]any{string(key), optionalID(data.ActorID), json.RawMessage(projected)})
	digest := sha256.Sum256(semantic)
	return entity.Activity{AuditID: auditID, SourceEventID: event.EventID, TenantID: tenant, AssetType: assetType, AssetID: event.AggregateID, AuditType: expectedType, Operation: data.Operation, Source: "identity", AssetVersion: data.New.Version, ActorID: optionalID(data.ActorID), PersonID: optionalID(data.New.PersonGID), Data: projected, SemanticHash: hex.EncodeToString(digest[:]), OccurredAt: event.OccurredAt.UTC()}, nil
}

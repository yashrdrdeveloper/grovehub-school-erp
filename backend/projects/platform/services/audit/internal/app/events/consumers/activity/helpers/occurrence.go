package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/google/uuid"

	client "grovehub/services/audit/client/producer"
	"grovehub/services/audit/internal/app/entity"
	validator "grovehub/services/audit/internal/app/validators/events"
)

func projectOccurrence(raw []byte) (entity.Activity, error) {
	event, err := validator.Occurrence(raw)
	if err != nil {
		return entity.Activity{}, err
	}

	data := event.Data
	identity := any(data.OccurrenceID)
	occurrence := &data.OccurrenceID

	if data.AssetVersion > 0 {
		if data.OccurrenceID != VersionEventID(data) {
			return entity.Activity{}, ErrInvalidEvent
		}
		identity = data.AssetVersion
		occurrence = nil
	}

	key, _ := json.Marshal([]any{data.TenantID, data.Source, data.AssetType, data.AssetID, client.ActivityOccurred, data.Operation, identity})
	id := uuid.NewSHA1(namespace, key).String()
	semantic, _ := json.Marshal([]any{string(key), data.ActorID, data.PersonID, data.OccurredAt.UTC()})
	digest := sha256.Sum256(semantic)

	var subject *string
	if data.AssetType == "ACCOUNT" {
		subject = &data.AssetID
	}

	return entity.Activity{
		AuditID: id,
		SourceEventID: event.EventID,
		TenantID: data.TenantID,
		AssetType: data.AssetType,
		AssetID: data.AssetID,
		AuditType: client.ActivityOccurred,
		Operation: data.Operation,
		Source: data.Source,
		OccurrenceID: occurrence,
		AssetVersion: data.AssetVersion,
		ActorID: optionalID(data.ActorID),
		PersonID: optionalID(data.PersonID),
		SubjectAccountID: subject,
		Data: json.RawMessage("{}"),
		SemanticHash: hex.EncodeToString(digest[:]),
		OccurredAt: data.OccurredAt.UTC(),
	}, nil
}

// VersionEventID derives an independent stable outbox ID, not a random UUID.
func VersionEventID(data client.Activity) string {
	key, _ := json.Marshal([]any{"versioned-source", data.TenantID, data.Source, data.AssetType, data.AssetID, data.Operation, data.AssetVersion})
	return uuid.NewSHA1(namespace, key).String()
}

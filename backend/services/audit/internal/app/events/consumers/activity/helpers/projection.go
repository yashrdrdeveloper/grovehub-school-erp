package helpers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
	coreevents "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events"
	"grovehub/services/audit/internal/app/entity"
)

var ErrInvalidEvent = errors.New("invalid or unscoped audit event")
var namespace = uuid.MustParse("bca86488-355c-4c40-a2cf-41121835ed83")

// A narrow schema, not a recursive blacklist. Unknown fields (including
// approval payload hashes, emails, credentials and free-form provider data)
// are never persisted. Existing role descriptions are deliberately omitted.
type snapshot struct {
	ID string `json:"id,omitempty"`
	AccountID string `json:"accountId,omitempty"`
	PersonGID string `json:"personGid,omitempty"`
	TenantID string `json:"tenantId,omitempty"`
	Version int64 `json:"version"`
	Status string `json:"status,omitempty"`
	Username *string `json:"username,omitempty"`
	CurrentStep int `json:"currentStep,omitempty"`
	StepCount int `json:"stepCount,omitempty"`
	Role *roleSnapshot `json:"role,omitempty"`
	Account *accountSnapshot `json:"account,omitempty"`
}

type roleSnapshot struct {
	Name string `json:"name"`
	Status string `json:"status"`
	Personas []string `json:"personas"`
	Permissions []string `json:"permissionCodes"`
}

type accountSnapshot struct {
	Status string `json:"status"`
	PrimaryRoleID string `json:"primaryRoleId"`
	SecondaryRoleIDs []string `json:"secondaryRoleIds"`
}

type payload struct {
	TenantID string `json:"tenantId"`
	ActorID string `json:"actorId"`
	Operation string `json:"operation"`
	RequestID string `json:"requestId"`
	RequestType string `json:"requestType"`
	Old *snapshot `json:"old"`
	New *snapshot `json:"new"`
}

func validID(value string) bool {
	id, err := uuid.Parse(value)
	return err == nil && id != uuid.Nil && id.String() == value
}

func optionalID(value string) *string {
	if !validID(value) { return nil }
	return &value
}

func Project(_ context.Context, expectedType string, raw []byte) (entity.Activity, error) {
	var event coreevents.Envelope[payload]
	if len(raw) > 512*1024 || json.Unmarshal(raw, &event) != nil || event.EventType != expectedType || event.Version != 1 || !validID(event.EventID) || !validID(event.AggregateID) || event.OccurredAt.IsZero() || event.Data.New == nil || event.Data.New.Version < 1 {
		return entity.Activity{}, ErrInvalidEvent
	}
	data := event.Data
	tenant := data.TenantID
	if tenant == "" { tenant = data.New.TenantID }
	if tenant == "" || strings.TrimSpace(tenant) != tenant || (data.New.TenantID != "" && data.New.TenantID != tenant) || (data.Old != nil && data.Old.TenantID != "" && data.Old.TenantID != tenant) || data.Operation == "" || len(data.Operation) > 80 {
		return entity.Activity{}, ErrInvalidEvent
	}
	assetType := ""
	switch expectedType {
	case "identity.account.changed.v1", "identity.account.registration.changed.v1", "identity.account.metadata.changed.v1", "identity.account.authority.changed.v1":
		assetType = "ACCOUNT"
	case "identity.role.changed.v1", "identity.role.metadata.changed.v1":
		assetType = "ROLE"
	case "identity.approval.changed.v1":
		assetType = "APPROVAL"
	default:
		return entity.Activity{}, ErrInvalidEvent
	}
	for _, value := range []*snapshot{data.Old,data.New} {
		if value == nil { continue }
		if (value.ID != "" && value.ID != event.AggregateID) || (value.AccountID != "" && value.AccountID != event.AggregateID) { return entity.Activity{}, ErrInvalidEvent }
		value.TenantID = ""
	}
	if assetType == "APPROVAL" && data.RequestID != event.AggregateID { return entity.Activity{}, ErrInvalidEvent }
	projected, err := json.Marshal(struct {
		Old *snapshot `json:"old"`
		New *snapshot `json:"new"`
		RequestID *string `json:"requestId,omitempty"`
		RequestType string `json:"requestType,omitempty"`
	}{data.Old, data.New, optionalID(data.RequestID), data.RequestType})
	if err != nil { return entity.Activity{}, ErrInvalidEvent }
	key, _ := json.Marshal([]any{tenant, assetType, event.AggregateID, expectedType, data.Operation, data.New.Version})
	auditID := uuid.NewSHA1(namespace, key).String()
	semantic, _ := json.Marshal([]any{string(key), optionalID(data.ActorID), json.RawMessage(projected)})
	digest := sha256.Sum256(semantic)
	return entity.Activity{AuditID: auditID, SourceEventID: event.EventID, TenantID: tenant, AssetType: assetType, AssetID: event.AggregateID, AuditType: expectedType, Operation: data.Operation, Source: "identity", AssetVersion: data.New.Version, ActorID: optionalID(data.ActorID), PersonID: optionalID(data.New.PersonGID), Data: projected, SemanticHash: hex.EncodeToString(digest[:]), OccurredAt: event.OccurredAt.UTC()}, nil
}

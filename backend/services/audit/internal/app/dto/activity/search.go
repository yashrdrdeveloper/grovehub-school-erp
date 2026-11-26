package activity

import (
	"time"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/request"
	"grovehub/services/audit/internal/app/entity"
)

type Cursor struct {
	OccurredAt time.Time `json:"occurredAt"`
	AuditID string `json:"auditId"`
}

type Search struct {
	AssetType string `json:"assetType,omitempty"`
	AssetID string `json:"assetId,omitempty"`
	ActorID string `json:"actorId,omitempty"`
	PersonID string `json:"personId,omitempty"`
	AuditType string `json:"auditType,omitempty"`
	Operation string `json:"operation,omitempty"`
	Source string `json:"source,omitempty"`
	Scope string `json:"scope,omitempty"`
	From *time.Time `json:"from,omitempty"`
	To *time.Time `json:"to,omitempty"`
	Limit int `json:"limit,omitempty"`
	Cursor *Cursor `json:"cursor,omitempty"`
}

type SearchRequest struct {
	request.Request
	RequestBody Search `validate:"valid"`
}

type Page struct {
	Items []entity.Activity `json:"items"`
	NextCursor *Cursor `json:"nextCursor,omitempty"`
}

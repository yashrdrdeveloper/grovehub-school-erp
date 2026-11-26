package audit

import (
	"time"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/request"
)

type Cursor struct {
	OccurredAt time.Time `json:"occurredAt"`
	Source string `json:"source"`
	EventID string `json:"eventId"`
}

type Search struct {
	Source string `json:"source,omitempty"`
	ResourceType string `json:"resourceType,omitempty"`
	ResourceID string `json:"resourceId,omitempty"`
	ActorID string `json:"actorId,omitempty"`
	Operation string `json:"operation,omitempty"`
	From *time.Time `json:"from,omitempty"`
	To *time.Time `json:"to,omitempty"`
	Limit int `json:"limit,omitempty"`
	Cursor *Cursor `json:"cursor,omitempty"`
}

type SearchRequest struct {
	request.Request
	RequestBody Search `validate:"valid"`
}

// Event intentionally excludes raw payloads, comments, credentials, hashes
// and request idempotency keys from this tenant-wide summary.
type Event struct {
	EventID string `json:"eventId"`
	Source string `json:"source"`
	ActorID string `json:"actorId"`
	ResourceType string `json:"resourceType"`
	ResourceID string `json:"resourceId"`
	RequestID *string `json:"requestId,omitempty"`
	Version *int64 `json:"version,omitempty"`
	Operation string `json:"operation"`
	OccurredAt time.Time `json:"occurredAt"`
}

type Page struct {
	Items []Event `json:"items"`
	NextCursor *Cursor `json:"nextCursor,omitempty"`
}

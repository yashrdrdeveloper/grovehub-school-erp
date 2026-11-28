package producer

import "time"

const ActivityOccurred = "audit.activity.occurred.v1"

// Activity describes a registered auditable operation. For versioned assets,
// supply AssetVersion and omit OccurrenceID; Audit derives the event identity.
// Otherwise generate OccurrenceID once per action and reuse it on retries.
// The registered operation determines which metadata may be persisted.
type Activity struct {
	OccurrenceID string `json:"occurrenceId"`
	TenantID string `json:"tenantId"`
	Source string `json:"source"`
	AssetType string `json:"assetType"`
	AssetID string `json:"assetId"`
	ActorID string `json:"actorId,omitempty"`
	PersonID string `json:"personId,omitempty"`
	AssetVersion int64 `json:"assetVersion,omitempty"`
	Operation string `json:"operation"`
	OccurredAt time.Time `json:"occurredAt"`
}

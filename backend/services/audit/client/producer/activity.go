package producer

import "time"

const ActivityOccurred = "audit.activity.occurred.v1"

// Activity describes an individual action without an asset revision.
// OccurrenceID must be generated once per action and reused on retries.
// The registered operation determines which metadata may be persisted.
type Activity struct {
	OccurrenceID string `json:"occurrenceId"`
	TenantID string `json:"tenantId"`
	Source string `json:"source"`
	AssetType string `json:"assetType"`
	AssetID string `json:"assetId"`
	ActorID string `json:"actorId"`
	Operation string `json:"operation"`
	OccurredAt time.Time `json:"occurredAt"`
}

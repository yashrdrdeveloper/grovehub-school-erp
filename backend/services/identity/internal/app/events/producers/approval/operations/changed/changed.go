package changed

import (
	"time"

	kind "grovehub/services/identity/internal/app/constants/approval"
)

type Snapshot struct {
	Status kind.Status `json:"status"`
	Version int64 `json:"version"`
	CurrentStep int `json:"currentStep"`
	StepCount int `json:"stepCount"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// Changed excludes raw intent, free-text reasons/comments, idempotency keys
// and session identifiers. Resource owners retain typed details in Identity.
type Changed struct {
	TenantID string `json:"tenantId"`
	RequestID string `json:"requestId"`
	RequestType kind.RequestType `json:"requestType"`
	PolicyID string `json:"policyId"`
	RequesterID string `json:"requesterId"`
	ActorID string `json:"actorId"`
	PayloadSHA256 string `json:"payloadSha256"`
	Operation string `json:"operation"`
	Old *Snapshot `json:"old"`
	New Snapshot `json:"new"`
}

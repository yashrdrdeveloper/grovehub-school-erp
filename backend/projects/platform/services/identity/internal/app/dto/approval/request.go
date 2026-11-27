package approval

import (
	"grovehub/services/identity/internal/app/constants/approval"
)

// LifecycleIntent is the first supported typed owner contract. The target
// version is part of the reviewed payload and cannot be changed by a reviewer.
type LifecycleIntent struct {
	TargetID string `json:"targetId"`
	Version int64 `json:"version"`
	Action string `json:"action"`
}

type CreateLifecycle struct {
	Type approval.RequestType `json:"type"`
	Intent LifecycleIntent `json:"intent"`
	Reason string `json:"reason"`
	IdempotencyKey string `json:"-"`
}

type Review struct {
	RequestID string `json:"requestId"`
	Version int64 `json:"version"`
	Decision approval.Decision `json:"decision"`
	Comment string `json:"comment"`
}

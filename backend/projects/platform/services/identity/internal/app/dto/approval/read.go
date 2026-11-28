package approval

import (
	"encoding/json"
	"time"

	kind "grovehub/services/identity/internal/app/constants/approval"
)

const (
	FieldRequestID = "requestId"
	FieldRequestType = "requestType"
	FieldStatus = "status"
	FieldCreatedAt = "createdAt"
	FieldUpdatedAt = "updatedAt"
	FieldExpiresAt = "expiresAt"
)

type Summary struct {
	RequestID string `json:"requestId"`
	RequesterID string `json:"requesterId"`
	RequestType kind.RequestType `json:"requestType"`
	Status kind.Status `json:"status"`
	CurrentStep int `json:"currentStep"`
	StepCount int `json:"stepCount"`
	Version int64 `json:"version"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type Detail struct {
	Summary
	PolicyID string `json:"policyId"`
	Intent json.RawMessage `json:"intent"`
	Reason string `json:"reason"`
}

type AuditEvent struct {
	EventID string `json:"eventId"`
	OccurredAt time.Time `json:"occurredAt"`
	Payload json.RawMessage `json:"payload"`
}

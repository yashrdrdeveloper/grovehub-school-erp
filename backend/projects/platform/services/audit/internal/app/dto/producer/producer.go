// Package producer owns the event payloads published through Audit.
package producer

import "time"

// Event contains a versioned domain envelope. Audit validates the schema and
// owns subject selection and transactional outbox publication.
type Event struct {
	Type    string
	Payload []byte
}

// Activity describes a registered auditable operation. For versioned assets,
// supply AssetVersion and omit OccurrenceID; Audit derives the event identity.
// Otherwise generate OccurrenceID once per action and reuse it on retries.
// The registered operation determines which metadata may be persisted.
type Activity struct {
	OccurrenceID string    `json:"occurrenceId"`
	TenantID     string    `json:"tenantId"`
	Source       string    `json:"source"`
	AssetType    string    `json:"assetType"`
	AssetID      string    `json:"assetId"`
	ActorID      string    `json:"actorId,omitempty"`
	PersonID     string    `json:"personId,omitempty"`
	AssetVersion int64     `json:"assetVersion,omitempty"`
	Operation    string    `json:"operation"`
	OccurredAt   time.Time `json:"occurredAt"`
}

// InvitationEvent contains only delivery lifecycle metadata, never mail content.
type InvitationEvent struct {
	TenantID      string     `json:"tenantId"`
	AccountID     string     `json:"accountId"`
	Purpose       string     `json:"purpose"`
	Operation     string     `json:"operation"`
	Attempt       int        `json:"attempt"`
	Reason        string     `json:"reason,omitempty"`
	NextAttemptAt *time.Time `json:"nextAttemptAt,omitempty"`
}

// OperatorActivity is an occurrence asserted by the independently authenticated
// Internal Tool workload. It is not a School identity or tenant activity.
type OperatorActivity struct {
	EventID       string    `json:"eventId"`
	ActorID       string    `json:"actorId"`
	ResourceID    string    `json:"resourceId"`
	Operation     string    `json:"operation"`
	CorrelationID *string   `json:"correlationId,omitempty"`
	OccurredAt    time.Time `json:"occurredAt"`
}

// SchoolRenameReceipt records a committed mutation, not an asserted user activity.
type SchoolRenameReceipt struct {
	ExecutionID     string    `json:"executionId"`
	SchoolID        string    `json:"schoolId"`
	PreviousVersion int64     `json:"previousVersion"`
	CorrelationID   *string   `json:"correlationId,omitempty"`
	OccurredAt      time.Time `json:"occurredAt"`
}

type SchoolStatusReceipt struct {
	ExecutionID     string    `json:"executionId"`
	SchoolID        string    `json:"schoolId"`
	PreviousVersion int64     `json:"previousVersion"`
	CorrelationID   *string   `json:"correlationId,omitempty"`
	OccurredAt      time.Time `json:"occurredAt"`
	Operation       string    `json:"operation"`
}

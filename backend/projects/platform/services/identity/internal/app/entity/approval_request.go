package entity

import (
	"encoding/json"
	"time"

	"grovehub/services/identity/internal/app/constants/approval"
)

// @Entity identity.approval_requests
type ApprovalRequest struct {
	RequestID string `gorm:"column:request_id;primaryKey;type:uuid"`
	TenantID string `gorm:"column:tenant_id"`
	PolicyID string `gorm:"column:policy_id;type:uuid"`
	RequesterID string `gorm:"column:requester_id;type:uuid"`
	RequestType approval.RequestType `gorm:"column:request_type"`
	Payload json.RawMessage `gorm:"column:payload;type:jsonb"`
	PayloadSHA256 string `gorm:"column:payload_sha256"`
	Reason string `gorm:"column:reason"`
	IdempotencyKey string `gorm:"column:idempotency_key"`
	Status approval.Status `gorm:"column:status"`
	CurrentStep int `gorm:"column:current_step"`
	StepCount int `gorm:"column:step_count"`
	Version int64 `gorm:"column:version"`
	ExpiresAt time.Time `gorm:"column:expires_at"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

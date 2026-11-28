package entity

import (
	"encoding/json"
	"time"
)

// @Entity audit.activities
type Activity struct {
	AuditID          string          `gorm:"column:audit_id;primaryKey;type:uuid" json:"auditId"`
	SourceEventID    string          `gorm:"column:source_event_id;type:uuid" json:"sourceEventId"`
	TenantID         string          `json:"-"`
	AssetType        string          `json:"assetType"`
	AssetID          string          `gorm:"type:uuid" json:"assetId"`
	AuditType        string          `json:"auditType"`
	Operation        string          `json:"operation"`
	Source           string          `json:"source"`
	OccurrenceID     *string         `gorm:"type:uuid" json:"occurrenceId,omitempty"`
	AssetVersion     int64           `json:"assetVersion"`
	ActorID          *string         `gorm:"type:uuid" json:"actorId,omitempty"`
	SubjectAccountID *string         `gorm:"type:uuid" json:"subjectAccountId,omitempty"`
	PersonID         *string         `gorm:"type:uuid" json:"personId,omitempty"`
	Data             json.RawMessage `gorm:"type:jsonb" json:"data"`
	SemanticHash     string          `json:"-"`
	OccurredAt       time.Time       `json:"occurredAt"`
	IngestedAt       time.Time       `json:"ingestedAt"`
}

package entity

import (
	"time"
	"grovehub/services/registry/internal/app/constants/status"
)

// @Entity registry.workloads
type Workload struct {
	WorkloadID string        `gorm:"column:workload_id;primaryKey;type:uuid"`
	ProjectID  *string       `gorm:"column:project_id;type:uuid"`
	Subject    string        `gorm:"column:subject;not null"`
	Issuer     string        `gorm:"column:issuer;not null"`
	Status     status.Status `gorm:"column:status;not null"`
	Version    int64         `gorm:"column:version;not null"`
	CreatedAt  time.Time     `gorm:"column:created_at;not null"`
	UpdatedAt  time.Time     `gorm:"column:updated_at;not null"`
}

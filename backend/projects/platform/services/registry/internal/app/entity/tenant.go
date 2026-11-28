package entity

import (
	"time"
	"grovehub/services/registry/internal/app/constants/status"
)

// @Entity registry.tenants
type Tenant struct {
	TenantID          string        `gorm:"column:tenant_id;primaryKey;type:uuid"`
	ProjectID         string        `gorm:"column:project_id;type:uuid;not null"`
	ExternalReference string        `gorm:"column:external_reference;not null"`
	Status            status.Status `gorm:"column:status;not null"`
	SessionGeneration int64         `gorm:"column:session_generation;not null"`
	Version           int64         `gorm:"column:version;not null"`
	CreatedAt         time.Time     `gorm:"column:created_at;not null"`
	UpdatedAt         time.Time     `gorm:"column:updated_at;not null"`
}

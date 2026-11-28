package entity

import (
	"time"
	"grovehub/libs/platform/scope"
)

// @Entity registry.workload_grants
type WorkloadGrant struct {
	GrantID         string      `gorm:"column:grant_id;primaryKey;type:uuid"`
	WorkloadID      string      `gorm:"column:workload_id;type:uuid;not null"`
	TargetProjectID *string     `gorm:"column:target_project_id;type:uuid"`
	TargetTenantID  *string     `gorm:"column:target_tenant_id;type:uuid"`
	Scope           scope.Level `gorm:"column:scope;not null"`
	Capability      string      `gorm:"column:capability;not null"`
	AllowDelegation bool        `gorm:"column:allow_delegation;not null"`
	CreatedAt       time.Time   `gorm:"column:created_at;not null"`
}

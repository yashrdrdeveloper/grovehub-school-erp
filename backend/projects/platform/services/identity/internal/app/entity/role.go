package entity

import "time"

type RoleStatus string

const (
	RoleActive   RoleStatus = "ACTIVE"
	RoleInactive RoleStatus = "INACTIVE"
)

// @Entity identity.roles
type Role struct {
	RoleID         string     `gorm:"column:role_id;primaryKey;type:uuid"`
	TenantID       string     `gorm:"column:tenant_id;not null"`
	Name           string     `gorm:"column:name;not null"`
	NormalizedName string     `gorm:"column:normalized_name;not null"`
	Description    string     `gorm:"column:description;not null"`
	Status         RoleStatus `gorm:"column:status;not null"`
	SearchDetails  string     `gorm:"column:search_details;not null"`
	Version        int64      `gorm:"column:version;not null"`
	CreatedBy      *string    `gorm:"column:created_by;type:uuid"`
	UpdatedBy      *string    `gorm:"column:updated_by;type:uuid"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at"`
}

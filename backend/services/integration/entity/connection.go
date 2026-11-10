package entity

import (
	connectionconstants "grovehub/services/integration/constants/connections"
	"time"
)

const (
	StatusPendingSetup = connectionconstants.PendingSetup
	StatusActive       = connectionconstants.Active
	StatusDisconnected = connectionconstants.Disconnected
)

// @Entity integration.connections
type Connection struct {
	ID             string     `gorm:"column:id;type:uuid;primaryKey"`
	TenantID       string     `gorm:"column:tenant_id;type:uuid;not null;index"`
	OwnerGID       string     `gorm:"column:owner_gid;type:uuid;not null;index"`
	Provider       string     `gorm:"column:provider;not null"`
	ExternalID     *string    `gorm:"column:external_id"`
	CredentialRef  *string    `gorm:"column:credential_ref"`
	Status         string     `gorm:"column:status;not null"`
	LastTestedAt   *time.Time `gorm:"column:last_tested_at"`
	LastTestStatus *string    `gorm:"column:last_test_status"`
	CreatedAt      time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;not null"`
}


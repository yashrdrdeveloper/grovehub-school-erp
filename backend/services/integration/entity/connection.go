package entity

import "time"

const (
 StatusPendingSetup = "PENDING_SETUP"
 StatusActive = "ACTIVE"
 StatusDisconnected = "DISCONNECTED"
)

type Connection struct {
 ID string `gorm:"column:id;type:uuid;primaryKey"`
 TenantID string `gorm:"column:tenant_id;type:uuid;not null;index"`
 OwnerGID string `gorm:"column:owner_gid;type:uuid;not null;index"`
 Provider string `gorm:"column:provider;not null"`
 CredentialRef *string `gorm:"column:credential_ref"`
 Status string `gorm:"column:status;not null"`
 CreatedAt time.Time `gorm:"column:created_at;not null"`
 UpdatedAt time.Time `gorm:"column:updated_at;not null"`
}
func (Connection) TableName() string { return "integration.connections" }

package entity

import "time"

// One cursor per provider resource and mailbox, not one cursor for all resources.
type SyncState struct {
 ID string `gorm:"column:id;type:uuid;primaryKey"`
 ConnectionID string `gorm:"column:connection_id;type:uuid;not null"`
 Resource string `gorm:"column:resource;not null"`
 Cursor string `gorm:"column:cursor;not null"`
 LastSyncedAt *time.Time `gorm:"column:last_synced_at"`
}
func (SyncState) TableName() string { return "integration.sync_states" }

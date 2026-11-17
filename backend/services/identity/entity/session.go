package entity

import "time"

// Session persists hashes of opaque bearer credentials; raw tokens leave only
// in login responses and never appear in persistence or application logs.
// @Entity identity.sessions
type Session struct {
    SessionID string `gorm:"column:session_id;type:uuid;primaryKey"`
    AccountId AccountId `gorm:"column:account_id;type:uuid;not null"`
    AccessHash string `gorm:"column:access_hash;type:char(64);not null"`
    RefreshHash string `gorm:"column:refresh_hash;type:char(64);not null"`
    AccessExpiresAt time.Time `gorm:"column:access_expires_at;not null"`
    RefreshExpiresAt time.Time `gorm:"column:refresh_expires_at;not null"`
    RevokedAt *time.Time `gorm:"column:revoked_at"`
    CreatedAt time.Time `gorm:"column:created_at;not null"`
    UpdatedAt time.Time `gorm:"column:updated_at;not null"`
}

func (Session) TableName() string { return "identity.sessions" }

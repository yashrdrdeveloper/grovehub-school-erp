package entity

import "time"

// @Entity identity.password_credentials
type PasswordCredential struct {
	AccountId         AccountId `gorm:"column:account_id;type:uuid;primaryKey"`
	PasswordHash      string    `gorm:"column:password_hash;not null"`
	PasswordChangedAt time.Time `gorm:"column:password_changed_at;not null"`
	Version           int64     `gorm:"column:version;not null"`
	CreatedAt         time.Time `gorm:"column:created_at;not null"`
	UpdatedAt         time.Time `gorm:"column:updated_at;not null"`
}

package entity

import "time"

type ContactType string

const (
	ContactTypeEmail ContactType = "EMAIL"
	ContactTypePhone ContactType = "PHONE"
)

// @Entity foundation.person_contacts
type Contact struct {
	PersonGID       string      `gorm:"column:person_gid;type:uuid;primaryKey"`
	Type            ContactType `gorm:"column:type;primaryKey"`
	Value           string      `gorm:"column:value;not null"`
	NormalizedValue string      `gorm:"column:normalized_value;primaryKey;not null"`
	IsPrimary       bool        `gorm:"column:is_primary;not null"`
	IsVerified      bool        `gorm:"column:is_verified;not null"`
	CreatedAt       time.Time   `gorm:"column:created_at;not null"`
	UpdatedAt       time.Time   `gorm:"column:updated_at;not null"`
}

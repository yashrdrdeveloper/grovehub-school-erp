package entity

import (
	"time"

	coreentity "grovehub/libs/core/entity"
)

type ContactType string

const (
	ContactTypeEmail ContactType = "EMAIL"
	ContactTypePhone ContactType = "PHONE"
)

// @Builder
// @Getter
// @Setter
// @GroveTable foundation.person_contacts
type Contact struct {
	PersonGID       GID         `gorm:"column:person_gid;type:uuid;primaryKey"`
	Type            ContactType `gorm:"column:type;primaryKey"`
	Value           string      `gorm:"column:value;not null"`
	NormalizedValue string      `gorm:"column:normalized_value;primaryKey;not null"`
	IsPrimary       bool        `gorm:"column:is_primary;not null"`
	IsVerified      bool        `gorm:"column:is_verified;not null"`
	CreatedAt       time.Time   `gorm:"column:created_at;not null"`
	UpdatedAt       time.Time   `gorm:"column:updated_at;not null"`
}

// ForPerson is intentional Foundation behavior, not generated boilerplate.
func (contact Contact) ForPerson(gid GID, createdAt time.Time, updatedAt time.Time) Contact {
	return coreentity.From(contact).With(func(value *Contact) {
		value.PersonGID = gid
		if value.CreatedAt.IsZero() {
			value.CreatedAt = createdAt
		}
		value.UpdatedAt = updatedAt
	}).Build()
}

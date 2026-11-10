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

func (Contact) TableName() string {
	return "foundation.person_contacts"
}

// NewContact preserves the current constructor contract; new use cases may use coreentity.New[Contact] directly.
func NewContact(contactType ContactType, value string, normalizedValue string, isPrimary bool, isVerified bool) Contact {
 return coreentity.New[Contact]().With(func(contact *Contact) {
  contact.Type = contactType
  contact.Value = value
  contact.NormalizedValue = normalizedValue
  contact.IsPrimary = isPrimary
  contact.IsVerified = isVerified
 }).Build()
}

func (contact Contact) ForPerson(gid GID, createdAt time.Time, updatedAt time.Time) Contact {
 return coreentity.From(contact).With(func(value *Contact) {
  value.PersonGID = gid
  if value.CreatedAt.IsZero() {
   value.CreatedAt = createdAt
  }
  value.UpdatedAt = updatedAt
 }).Build()
}

func (contact Contact) GetPersonGID() GID {
	return contact.PersonGID
}

func (contact Contact) GetType() ContactType {
	return contact.Type
}

func (contact Contact) GetValue() string {
	return contact.Value
}

func (contact Contact) GetNormalizedValue() string {
	return contact.NormalizedValue
}

func (contact Contact) GetIsPrimary() bool {
	return contact.IsPrimary
}

func (contact Contact) GetIsVerified() bool {
	return contact.IsVerified
}

func (contact Contact) GetCreatedAt() time.Time {
	return contact.CreatedAt
}

func (contact Contact) GetUpdatedAt() time.Time {
	return contact.UpdatedAt
}

package entity

import "time"

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

func NewContact(
	contactType ContactType,
	value string,
	normalizedValue string,
	isPrimary bool,
	isVerified bool,
) Contact {
	return Contact{
		Type:            contactType,
		Value:           value,
		NormalizedValue: normalizedValue,
		IsPrimary:       isPrimary,
		IsVerified:      isVerified,
	}
}

func (contact Contact) ForPerson(gid GID, createdAt time.Time, updatedAt time.Time) Contact {
	contact.PersonGID = gid
	if contact.CreatedAt.IsZero() {
		contact.CreatedAt = createdAt
	}
	contact.UpdatedAt = updatedAt
	return contact
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

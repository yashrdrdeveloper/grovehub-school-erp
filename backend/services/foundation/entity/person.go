package entity

import "time"

type PersonStatus string

const (
	PersonStatusActive      PersonStatus = "ACTIVE"
	PersonStatusDeactivated PersonStatus = "DEACTIVATED"
)

type Person struct {
	GID           GID          `gorm:"column:gid;type:uuid;primaryKey"`
	FirstName     string       `gorm:"column:first_name;not null"`
	MiddleName    *string      `gorm:"column:middle_name"`
	LastName      string       `gorm:"column:last_name;not null"`
	Contacts      []Contact    `gorm:"-"`
	Status        PersonStatus `gorm:"column:status;not null"`
	Version       int64        `gorm:"column:version;not null"`
	SearchDetails string       `gorm:"column:search_details;not null"`
	CreatedAt     time.Time    `gorm:"column:created_at;not null"`
	UpdatedAt     time.Time    `gorm:"column:updated_at;not null"`
}

func (Person) TableName() string {
	return "foundation.persons"
}

func NewPerson(
	gid GID,
	firstName string,
	middleName string,
	lastName string,
	contacts []Contact,
	status PersonStatus,
	version int64,
	createdAt time.Time,
	updatedAt time.Time,
) Person {
	return Person{
		GID:        gid,
		FirstName:  firstName,
		MiddleName: optionalString(middleName),
		LastName:   lastName,
		Contacts:   copyContacts(contacts),
		Status:     status,
		Version:    version,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}

func (person Person) GetGID() GID {
	return person.GID
}

func (person Person) GetFirstName() string {
	return person.FirstName
}

func (person Person) GetMiddleName() string {
	if person.MiddleName == nil {
		return ""
	}
	return *person.MiddleName
}

func (person Person) GetLastName() string {
	return person.LastName
}

func (person Person) GetContacts() []Contact {
	return copyContacts(person.Contacts)
}

func (person Person) GetStatus() PersonStatus {
	return person.Status
}

func (person Person) GetVersion() int64 {
	return person.Version
}

func (person Person) GetCreatedAt() time.Time {
	return person.CreatedAt
}

func (person Person) GetUpdatedAt() time.Time {
	return person.UpdatedAt
}

func (person Person) SetFirstName(firstName string) Person {
	person.FirstName = firstName
	return person
}

func (person Person) SetMiddleName(middleName string) Person {
	person.MiddleName = optionalString(middleName)
	return person
}

func (person Person) SetLastName(lastName string) Person {
	person.LastName = lastName
	return person
}

func (person Person) SetContacts(contacts []Contact) Person {
	person.Contacts = copyContacts(contacts)
	return person
}

func (person Person) SetStatus(status PersonStatus) Person {
	person.Status = status
	return person
}

func (person Person) SetVersion(version int64) Person {
	person.Version = version
	return person
}

func (person Person) SetUpdatedAt(updatedAt time.Time) Person {
	person.UpdatedAt = updatedAt
	return person
}

func (person Person) BuildContacts() []Contact {
	contacts := copyContacts(person.Contacts)
	for index := range contacts {
		contacts[index] = contacts[index].ForPerson(person.GID, person.UpdatedAt, person.UpdatedAt)
	}
	return contacts
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func copyContacts(contacts []Contact) []Contact {
	if len(contacts) == 0 {
		return nil
	}

	copied := make([]Contact, len(contacts))
	copy(copied, contacts)
	return copied
}

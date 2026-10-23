package entity

import "time"

type PersonStatus string

const (
	PersonStatusActive      PersonStatus = "ACTIVE"
	PersonStatusDeactivated PersonStatus = "DEACTIVATED"
)

type Person struct {
	gid        GID
	firstName  string
	middleName string
	lastName   string
	contacts   []Contact
	status     PersonStatus
	version    int64
	createdAt  time.Time
	updatedAt  time.Time
}

func NewPerson(gid GID, firstName string, middleName string, lastName string, contacts []Contact, status PersonStatus, version int64, createdAt time.Time, updatedAt time.Time) Person {
	return Person{gid: gid, firstName: firstName, middleName: middleName, lastName: lastName, contacts: copyContacts(contacts), status: status, version: version, createdAt: createdAt, updatedAt: updatedAt}
}

func (person Person) GetGID() GID { return person.gid }
func (person Person) GetFirstName() string { return person.firstName }
func (person Person) GetMiddleName() string { return person.middleName }
func (person Person) GetLastName() string { return person.lastName }
func (person Person) GetContacts() []Contact { return copyContacts(person.contacts) }
func (person Person) GetStatus() PersonStatus { return person.status }
func (person Person) GetVersion() int64 { return person.version }
func (person Person) GetCreatedAt() time.Time { return person.createdAt }
func (person Person) GetUpdatedAt() time.Time { return person.updatedAt }
func (person Person) SetFirstName(firstName string) Person { person.firstName = firstName; return person }
func (person Person) SetMiddleName(middleName string) Person { person.middleName = middleName; return person }
func (person Person) SetLastName(lastName string) Person { person.lastName = lastName; return person }
func (person Person) SetContacts(contacts []Contact) Person { person.contacts = copyContacts(contacts); return person }
func (person Person) SetStatus(status PersonStatus) Person { person.status = status; return person }
func (person Person) SetVersion(version int64) Person { person.version = version; return person }
func (person Person) SetUpdatedAt(updatedAt time.Time) Person { person.updatedAt = updatedAt; return person }

func copyContacts(contacts []Contact) []Contact {
	if len(contacts) == 0 { return nil }
	copied := make([]Contact, len(contacts))
	copy(copied, contacts)
	return copied
}

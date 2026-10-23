package entity

import "time"

type PersonStatus string

const (
	PersonStatusActive      PersonStatus = "ACTIVE"
	PersonStatusDeactivated PersonStatus = "DEACTIVATED"
)

type Person struct {
	gid       GID
	status    PersonStatus
	version   int64
	createdAt time.Time
	updatedAt time.Time
}

func NewPerson(
	gid GID,
	status PersonStatus,
	version int64,
	createdAt time.Time,
	updatedAt time.Time,
) Person {
	return Person{
		gid:       gid,
		status:    status,
		version:   version,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (person Person) GetGID() GID {
	return person.gid
}

func (person Person) GetStatus() PersonStatus {
	return person.status
}

func (person Person) GetVersion() int64 {
	return person.version
}

func (person Person) GetCreatedAt() time.Time {
	return person.createdAt
}

func (person Person) GetUpdatedAt() time.Time {
	return person.updatedAt
}

func (person Person) SetStatus(status PersonStatus) Person {
	person.status = status
	return person
}

func (person Person) SetVersion(version int64) Person {
	person.version = version
	return person
}

func (person Person) SetUpdatedAt(updatedAt time.Time) Person {
	person.updatedAt = updatedAt
	return person
}

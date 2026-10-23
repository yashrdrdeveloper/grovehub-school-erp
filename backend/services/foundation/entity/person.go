package entity

import (
	"errors"
	"time"
)

type PersonStatus string

const (
	PersonStatusActive      PersonStatus = "ACTIVE"
	PersonStatusDeactivated PersonStatus = "DEACTIVATED"
)

var (
	ErrPersonGIDRequired     = errors.New("person gid is required")
	ErrInvalidPersonStatus   = errors.New("invalid person status")
	ErrInvalidPersonVersion  = errors.New("person version must be at least 1")
	ErrPersonTimestampNeeded = errors.New("person timestamp is required")
)

type Person struct {
	GID       GID
	Status    PersonStatus
	Version   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPerson(gid GID, now time.Time) (Person, error) {
	if gid.IsZero() {
		return Person{}, ErrPersonGIDRequired
	}
	if now.IsZero() {
		return Person{}, ErrPersonTimestampNeeded
	}

	timestamp := now.UTC()

	return Person{
		GID:       gid,
		Status:    PersonStatusActive,
		Version:   1,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}, nil
}

func (status PersonStatus) IsValid() bool {
	return status == PersonStatusActive || status == PersonStatusDeactivated
}

func (person Person) Validate() error {
	if person.GID.IsZero() {
		return ErrPersonGIDRequired
	}
	if !person.Status.IsValid() {
		return ErrInvalidPersonStatus
	}
	if person.Version < 1 {
		return ErrInvalidPersonVersion
	}
	if person.CreatedAt.IsZero() || person.UpdatedAt.IsZero() {
		return ErrPersonTimestampNeeded
	}

	return nil
}

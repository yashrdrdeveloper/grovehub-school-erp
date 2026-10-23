package entity

import "time"

type PersonStatus string

const (
	PersonStatusActive      PersonStatus = "ACTIVE"
	PersonStatusDeactivated PersonStatus = "DEACTIVATED"
)

type Person struct {
	GID       GID
	Status    PersonStatus
	Version   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

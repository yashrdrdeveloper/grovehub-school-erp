package entity

import "time"

type PersonStatus string

const (
	PersonStatusActive      PersonStatus = "ACTIVE"
	PersonStatusDeactivated PersonStatus = "DEACTIVATED"
)

// @Entity foundation.persons
type Person struct {
	GID                  string       `gorm:"column:gid;type:uuid;primaryKey"`
	FirstName            string       `gorm:"column:first_name;not null"`
	MiddleName           *string      `gorm:"column:middle_name"`
	LastName             string       `gorm:"column:last_name;not null"`
	Status               PersonStatus `gorm:"column:status;not null"`
	Version              int64        `gorm:"column:version;not null"`
	SearchDetails        string       `gorm:"column:search_details;not null"`
	IdempotencyKey       *string      `gorm:"column:create_idempotency_key"`
	RequestHash          *string      `gorm:"column:create_request_hash"`
	CreatedAt            time.Time    `gorm:"column:created_at;not null"`
	UpdatedAt            time.Time    `gorm:"column:updated_at;not null"`
}

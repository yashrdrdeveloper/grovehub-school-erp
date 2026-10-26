package impl

import "time"

const (
	personTable        = "foundation.persons"
	personContactTable = "foundation.person_contacts"
)

type personModel struct {
	GID           string    `gorm:"column:gid;primaryKey"`
	FirstName     string    `gorm:"column:first_name"`
	MiddleName    *string   `gorm:"column:middle_name"`
	LastName      string    `gorm:"column:last_name"`
	Status        string    `gorm:"column:status"`
	Version       int64     `gorm:"column:version"`
	SearchDetails string    `gorm:"column:search_details"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (personModel) TableName() string {
	return personTable
}

type personContactModel struct {
	PersonGID       string    `gorm:"column:person_gid;primaryKey"`
	Type            string    `gorm:"column:type;primaryKey"`
	Value           string    `gorm:"column:value"`
	NormalizedValue string    `gorm:"column:normalized_value;primaryKey"`
	IsPrimary       bool      `gorm:"column:is_primary"`
	IsVerified      bool      `gorm:"column:is_verified"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

func (personContactModel) TableName() string {
	return personContactTable
}

package entity

import "time"

type AccountId string
type PersonGid string
type AccountStatus string

const (
	AccountStatusActive      AccountStatus = "ACTIVE"
	AccountStatusSuspended   AccountStatus = "SUSPENDED"
	AccountStatusDeactivated AccountStatus = "DEACTIVATED"
)

// @Entity identity.accounts
type Account struct {
	AccountId                  AccountId     `gorm:"column:account_id;type:uuid;primaryKey"`
	RegistrationIdempotencyKey *string       `gorm:"column:registration_idempotency_key"`
	RegistrationRequestHash    *string       `gorm:"column:registration_request_hash"`
	PersonGid                  PersonGid     `gorm:"column:person_gid;type:uuid;not null"`
	Email                      string        `gorm:"column:email;not null"`
	Phone                      *string       `gorm:"column:phone"`
	Status                     AccountStatus `gorm:"column:status;not null"`
	Version                    int64         `gorm:"column:version;not null"`
	CreatedAt                  time.Time     `gorm:"column:created_at;not null"`
	UpdatedAt                  time.Time     `gorm:"column:updated_at;not null"`
}

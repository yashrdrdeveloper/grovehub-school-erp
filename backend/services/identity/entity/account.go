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

type Account struct {
	AccountId AccountId     `gorm:"column:account_id;type:uuid;primaryKey"`
	PersonGid PersonGid     `gorm:"column:person_gid;type:uuid;not null"`
	Email     string        `gorm:"column:email;not null"`
	Phone     *string       `gorm:"column:phone"`
	Status    AccountStatus `gorm:"column:status;not null"`
	Version   int64         `gorm:"column:version;not null"`
	CreatedAt time.Time     `gorm:"column:created_at;not null"`
	UpdatedAt time.Time     `gorm:"column:updated_at;not null"`
}

func (Account) TableName() string {
	return "identity.accounts"
}

func NewAccount(
	accountId AccountId,
	personGid PersonGid,
	email string,
	phone string,
	status AccountStatus,
	version int64,
	createdAt time.Time,
	updatedAt time.Time,
) Account {
	return Account{
		AccountId: accountId,
		PersonGid: personGid,
		Email:     email,
		Phone:     optionalString(phone),
		Status:    status,
		Version:   version,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func (account Account) GetAccountId() AccountId { return account.AccountId }
func (account Account) GetPersonGid() PersonGid { return account.PersonGid }
func (account Account) GetEmail() string        { return account.Email }
func (account Account) GetPhone() string {
	if account.Phone == nil {
		return ""
	}
	return *account.Phone
}
func (account Account) GetStatus() AccountStatus { return account.Status }
func (account Account) GetVersion() int64        { return account.Version }
func (account Account) GetCreatedAt() time.Time  { return account.CreatedAt }
func (account Account) GetUpdatedAt() time.Time  { return account.UpdatedAt }

func (account Account) SetEmail(email string) Account {
	account.Email = email
	return account
}

func (account Account) SetPhone(phone string) Account {
	account.Phone = optionalString(phone)
	return account
}

func (account Account) SetStatus(status AccountStatus) Account {
	account.Status = status
	return account
}

func (account Account) SetVersion(version int64) Account {
	account.Version = version
	return account
}

func (account Account) SetUpdatedAt(updatedAt time.Time) Account {
	account.UpdatedAt = updatedAt
	return account
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

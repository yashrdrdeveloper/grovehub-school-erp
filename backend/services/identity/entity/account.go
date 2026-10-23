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
	accountId AccountId
	personGid PersonGid
	email     string
	phone     string
	status    AccountStatus
	version   int64
	createdAt time.Time
	updatedAt time.Time
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
		accountId: accountId,
		personGid: personGid,
		email:     email,
		phone:     phone,
		status:    status,
		version:   version,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (account Account) GetAccountId() AccountId {
	return account.accountId
}

func (account Account) GetPersonGid() PersonGid {
	return account.personGid
}

func (account Account) GetEmail() string {
	return account.email
}

func (account Account) GetPhone() string {
	return account.phone
}

func (account Account) GetStatus() AccountStatus {
	return account.status
}

func (account Account) GetVersion() int64 {
	return account.version
}

func (account Account) GetCreatedAt() time.Time {
	return account.createdAt
}

func (account Account) GetUpdatedAt() time.Time {
	return account.updatedAt
}

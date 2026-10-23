package entity

import "time"

type AccountId string
type PersonGid string
type AccountStatus string

const (
	AccountStatusPendingVerification AccountStatus = "PENDING_VERIFICATION"
	AccountStatusActive              AccountStatus = "ACTIVE"
	AccountStatusLocked              AccountStatus = "LOCKED"
	AccountStatusSuspended           AccountStatus = "SUSPENDED"
	AccountStatusDeactivated         AccountStatus = "DEACTIVATED"
)

type Account struct {
	accountId AccountId
	personGid PersonGid
	status    AccountStatus
	version   int64
	createdAt time.Time
	updatedAt time.Time
}

func NewAccount(
	accountId AccountId,
	personGid PersonGid,
	status AccountStatus,
	version int64,
	createdAt time.Time,
	updatedAt time.Time,
) Account {
	return Account{
		accountId: accountId,
		personGid: personGid,
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

func (account Account) SetStatus(status AccountStatus) Account {
	account.status = status
	return account
}

func (account Account) SetVersion(version int64) Account {
	account.version = version
	return account
}

func (account Account) SetUpdatedAt(updatedAt time.Time) Account {
	account.updatedAt = updatedAt
	return account
}

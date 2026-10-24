package entity

import "time"

type PasswordCredential struct {
	accountId         AccountId
	passwordHash      string
	passwordChangedAt time.Time
	version           int64
	createdAt         time.Time
	updatedAt         time.Time
}

func NewPasswordCredential(
	accountId AccountId,
	passwordHash string,
	passwordChangedAt time.Time,
	version int64,
	createdAt time.Time,
	updatedAt time.Time,
) PasswordCredential {
	return PasswordCredential{
		accountId:         accountId,
		passwordHash:      passwordHash,
		passwordChangedAt: passwordChangedAt,
		version:           version,
		createdAt:         createdAt,
		updatedAt:         updatedAt,
	}
}

func (credential PasswordCredential) GetAccountId() AccountId {
	return credential.accountId
}

func (credential PasswordCredential) GetPasswordHash() string {
	return credential.passwordHash
}

func (credential PasswordCredential) GetPasswordChangedAt() time.Time {
	return credential.passwordChangedAt
}

func (credential PasswordCredential) GetVersion() int64 {
	return credential.version
}

func (credential PasswordCredential) GetCreatedAt() time.Time {
	return credential.createdAt
}

func (credential PasswordCredential) GetUpdatedAt() time.Time {
	return credential.updatedAt
}

func (credential PasswordCredential) SetPasswordHash(passwordHash string) PasswordCredential {
	credential.passwordHash = passwordHash
	return credential
}

func (credential PasswordCredential) SetPasswordChangedAt(passwordChangedAt time.Time) PasswordCredential {
	credential.passwordChangedAt = passwordChangedAt
	return credential
}

func (credential PasswordCredential) SetVersion(version int64) PasswordCredential {
	credential.version = version
	return credential
}

func (credential PasswordCredential) SetUpdatedAt(updatedAt time.Time) PasswordCredential {
	credential.updatedAt = updatedAt
	return credential
}

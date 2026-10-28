package entity

import "time"

type PasswordCredential struct {
	AccountId         AccountId `gorm:"column:account_id;type:uuid;primaryKey"`
	PasswordHash      string    `gorm:"column:password_hash;not null"`
	PasswordChangedAt time.Time `gorm:"column:password_changed_at;not null"`
	Version           int64     `gorm:"column:version;not null"`
	CreatedAt         time.Time `gorm:"column:created_at;not null"`
	UpdatedAt         time.Time `gorm:"column:updated_at;not null"`
}

func (PasswordCredential) TableName() string {
	return "identity.password_credentials"
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
		AccountId:         accountId,
		PasswordHash:      passwordHash,
		PasswordChangedAt: passwordChangedAt,
		Version:           version,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}
}

func (credential PasswordCredential) GetAccountId() AccountId { return credential.AccountId }
func (credential PasswordCredential) GetPasswordHash() string { return credential.PasswordHash }
func (credential PasswordCredential) GetPasswordChangedAt() time.Time {
	return credential.PasswordChangedAt
}
func (credential PasswordCredential) GetVersion() int64       { return credential.Version }
func (credential PasswordCredential) GetCreatedAt() time.Time { return credential.CreatedAt }
func (credential PasswordCredential) GetUpdatedAt() time.Time { return credential.UpdatedAt }

func (credential PasswordCredential) SetPasswordHash(passwordHash string) PasswordCredential {
	credential.PasswordHash = passwordHash
	return credential
}

func (credential PasswordCredential) SetPasswordChangedAt(passwordChangedAt time.Time) PasswordCredential {
	credential.PasswordChangedAt = passwordChangedAt
	return credential
}

func (credential PasswordCredential) SetVersion(version int64) PasswordCredential {
	credential.Version = version
	return credential
}

func (credential PasswordCredential) SetUpdatedAt(updatedAt time.Time) PasswordCredential {
	credential.UpdatedAt = updatedAt
	return credential
}

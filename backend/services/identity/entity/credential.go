package entity

import "time"

type CredentialId string
type CredentialType string
type CredentialAlgorithm string
type CredentialStatus string

const (
	CredentialTypePassword CredentialType = "PASSWORD"
)

const (
	CredentialAlgorithmArgon2Id CredentialAlgorithm = "ARGON2ID"
)

const (
	CredentialStatusActive  CredentialStatus = "ACTIVE"
	CredentialStatusRevoked CredentialStatus = "REVOKED"
)

type Credential struct {
	credentialId   CredentialId
	accountId      AccountId
	credentialType CredentialType
	secretHash     string
	algorithm      CredentialAlgorithm
	status         CredentialStatus
	lastUsedAt     *time.Time
	createdAt      time.Time
	updatedAt      time.Time
}

func NewCredential(
	credentialId CredentialId,
	accountId AccountId,
	credentialType CredentialType,
	secretHash string,
	algorithm CredentialAlgorithm,
	status CredentialStatus,
	lastUsedAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) Credential {
	return Credential{
		credentialId:   credentialId,
		accountId:      accountId,
		credentialType: credentialType,
		secretHash:     secretHash,
		algorithm:      algorithm,
		status:         status,
		lastUsedAt:     copyTime(lastUsedAt),
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}

func (credential Credential) GetCredentialId() CredentialId {
	return credential.credentialId
}

func (credential Credential) GetAccountId() AccountId {
	return credential.accountId
}

func (credential Credential) GetType() CredentialType {
	return credential.credentialType
}

func (credential Credential) GetSecretHash() string {
	return credential.secretHash
}

func (credential Credential) GetAlgorithm() CredentialAlgorithm {
	return credential.algorithm
}

func (credential Credential) GetStatus() CredentialStatus {
	return credential.status
}

func (credential Credential) GetLastUsedAt() *time.Time {
	return copyTime(credential.lastUsedAt)
}

func (credential Credential) GetCreatedAt() time.Time {
	return credential.createdAt
}

func (credential Credential) GetUpdatedAt() time.Time {
	return credential.updatedAt
}

func (credential Credential) SetStatus(status CredentialStatus) Credential {
	credential.status = status
	return credential
}

func (credential Credential) SetLastUsedAt(lastUsedAt time.Time) Credential {
	credential.lastUsedAt = &lastUsedAt
	return credential
}

func (credential Credential) SetUpdatedAt(updatedAt time.Time) Credential {
	credential.updatedAt = updatedAt
	return credential
}

func copyTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}

	copied := *value
	return &copied
}

package entity

import "time"

type LoginIdentifierId string
type LoginIdentifierType string
type LoginIdentifierStatus string

const (
	LoginIdentifierTypeEmail LoginIdentifierType = "EMAIL"
	LoginIdentifierTypePhone LoginIdentifierType = "PHONE"
)

const (
	LoginIdentifierStatusActive   LoginIdentifierStatus = "ACTIVE"
	LoginIdentifierStatusDisabled LoginIdentifierStatus = "DISABLED"
)

type LoginIdentifier struct {
	identifierId    LoginIdentifierId
	accountId       AccountId
	identifierType  LoginIdentifierType
	value           string
	normalizedValue string
	isPrimary       bool
	isVerified      bool
	status          LoginIdentifierStatus
	createdAt       time.Time
	updatedAt       time.Time
}

func NewLoginIdentifier(
	identifierId LoginIdentifierId,
	accountId AccountId,
	identifierType LoginIdentifierType,
	value string,
	normalizedValue string,
	isPrimary bool,
	isVerified bool,
	status LoginIdentifierStatus,
	createdAt time.Time,
	updatedAt time.Time,
) LoginIdentifier {
	return LoginIdentifier{
		identifierId:    identifierId,
		accountId:       accountId,
		identifierType:  identifierType,
		value:           value,
		normalizedValue: normalizedValue,
		isPrimary:       isPrimary,
		isVerified:      isVerified,
		status:          status,
		createdAt:       createdAt,
		updatedAt:       updatedAt,
	}
}

func (identifier LoginIdentifier) GetIdentifierId() LoginIdentifierId {
	return identifier.identifierId
}

func (identifier LoginIdentifier) GetAccountId() AccountId {
	return identifier.accountId
}

func (identifier LoginIdentifier) GetType() LoginIdentifierType {
	return identifier.identifierType
}

func (identifier LoginIdentifier) GetValue() string {
	return identifier.value
}

func (identifier LoginIdentifier) GetNormalizedValue() string {
	return identifier.normalizedValue
}

func (identifier LoginIdentifier) GetIsPrimary() bool {
	return identifier.isPrimary
}

func (identifier LoginIdentifier) GetIsVerified() bool {
	return identifier.isVerified
}

func (identifier LoginIdentifier) GetStatus() LoginIdentifierStatus {
	return identifier.status
}

func (identifier LoginIdentifier) GetCreatedAt() time.Time {
	return identifier.createdAt
}

func (identifier LoginIdentifier) GetUpdatedAt() time.Time {
	return identifier.updatedAt
}

func (identifier LoginIdentifier) SetIsPrimary(isPrimary bool) LoginIdentifier {
	identifier.isPrimary = isPrimary
	return identifier
}

func (identifier LoginIdentifier) SetIsVerified(isVerified bool) LoginIdentifier {
	identifier.isVerified = isVerified
	return identifier
}

func (identifier LoginIdentifier) SetStatus(status LoginIdentifierStatus) LoginIdentifier {
	identifier.status = status
	return identifier
}

func (identifier LoginIdentifier) SetUpdatedAt(updatedAt time.Time) LoginIdentifier {
	identifier.updatedAt = updatedAt
	return identifier
}

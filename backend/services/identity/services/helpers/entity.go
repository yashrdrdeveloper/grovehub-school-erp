package helpers

import (
	"time"

	"grovehub/services/identity/entity"
	identityservices "grovehub/services/identity/services"
)

// OptionalString preserves NULL semantics for optional registration fields.
func OptionalString(value string) *string {
	if value == "" { return nil }
	return &value
}

// NullableString is the service-layer read counterpart.
func NullableString(value *string) string {
	if value == nil { return "" }
	return *value
}

// NewAccount owns account creation defaults and idempotency metadata.
func NewAccount(id entity.AccountId, input identityservices.CreateAccountInput, at time.Time) entity.Account {
	return entity.Account{
		AccountId: id,
		PersonGid: entity.PersonGid(input.PersonGid),
		Email: input.Email,
		Phone: OptionalString(input.Phone),
		RegistrationIdempotencyKey: OptionalString(input.RegistrationIdempotencyKey),
		RegistrationRequestHash: OptionalString(input.RegistrationRequestHash),
		Status: entity.AccountStatusActive,
		Version: 1,
		CreatedAt: at,
		UpdatedAt: at,
	}
}

// NewPasswordCredential creates a credential after the service has hashed the password.
func NewPasswordCredential(accountID entity.AccountId, passwordHash string, at time.Time) entity.PasswordCredential {
	return entity.PasswordCredential{
		AccountId: accountID,
		PasswordHash: passwordHash,
		PasswordChangedAt: at,
		Version: 1,
		CreatedAt: at,
		UpdatedAt: at,
	}
}

package impl

import (
	"grovehub/services/identity/entity"
	identityvalidators "grovehub/services/identity/validators"
	validatorhelpers "grovehub/services/identity/validators/helpers"
)

type PasswordCredentialValidator struct{}

func NewPasswordCredentialValidator() identityvalidators.PasswordCredentialValidator {
	return &PasswordCredentialValidator{}
}

func (validator *PasswordCredentialValidator) ValidateAccountId(value string) (entity.AccountId, error) {
	return validatorhelpers.NormalizeAccountId(value)
}

func (validator *PasswordCredentialValidator) ValidatePassword(value string) error {
	return validatorhelpers.ValidatePassword(value)
}

func (validator *PasswordCredentialValidator) ValidatePasswordCandidate(value string) error {
	return validatorhelpers.ValidatePasswordCandidate(value)
}

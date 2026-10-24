package validators

import "grovehub/services/identity/entity"

type PasswordCredentialValidator interface {
	ValidateAccountId(value string) (entity.AccountId, error)
	ValidatePassword(value string) error
	ValidatePasswordCandidate(value string) error
}

package validators

import (
	"grovehub/services/identity/entity"
	identityservices "grovehub/services/identity/services"
)

type ValidatedCreateAccount struct {
	PersonGid entity.PersonGid
	Email     string
	Phone     string
}

type ValidatedUpdateAccount struct {
	Email *string
	Phone *string
}

type AccountValidator interface {
	ValidateCreate(input identityservices.CreateAccountInput) (ValidatedCreateAccount, error)
	ValidateUpdate(input identityservices.UpdateAccountInput) (ValidatedUpdateAccount, error)
	ValidateAccountId(value string) (entity.AccountId, error)
	ValidatePersonGid(value string) (entity.PersonGid, error)
	ValidateEmail(value string) (string, error)
}

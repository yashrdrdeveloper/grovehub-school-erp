package validators

import (
	"grovehub/services/identity/entity"
	identityservices "grovehub/services/identity/services"
)

type AccountValidator interface {
	ValidateCreate(input identityservices.CreateAccountInput) (identityservices.CreateAccountInput, error)
	ValidateUpdate(input identityservices.UpdateAccountInput) (identityservices.UpdateAccountInput, error)
	ValidateAccountId(value string) (entity.AccountId, error)
	ValidatePersonGid(value string) (entity.PersonGid, error)
	ValidateEmail(value string) (string, error)
}

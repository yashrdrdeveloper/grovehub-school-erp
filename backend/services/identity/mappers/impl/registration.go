package impl

import (
	identityresponse "grovehub/services/identity/dto/response"
	"grovehub/services/identity/entity"
	identitymappers "grovehub/services/identity/mappers"
)

type RegistrationMapper struct{}

func NewRegistrationMapper() identitymappers.RegistrationMapper {
	return &RegistrationMapper{}
}

func (mapper *RegistrationMapper) ToRegistrationResponse(account entity.Account) identityresponse.Registration {
	return identityresponse.Registration{
		AccountId: string(account.GetAccountId()),
		PersonGid: string(account.GetPersonGid()),
		Email:     account.GetEmail(),
		Phone:     account.GetPhone(),
		Status:    string(account.GetStatus()),
	}
}

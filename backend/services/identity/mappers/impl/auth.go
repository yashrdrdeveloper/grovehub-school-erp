package impl

import (
	identityresponse "grovehub/services/identity/dto/response"
	"grovehub/services/identity/entity"
	identitymappers "grovehub/services/identity/mappers"
)

type AuthMapper struct{}

func NewAuthMapper() identitymappers.AuthMapper {
	return &AuthMapper{}
}

func (mapper *AuthMapper) ToRegisterResponse(account entity.Account) identityresponse.Register {
	return identityresponse.Register{
		AccountId: string(account.GetAccountId()),
		PersonGid: string(account.GetPersonGid()),
		Email:     account.GetEmail(),
		Phone:     account.GetPhone(),
		Status:    string(account.GetStatus()),
	}
}

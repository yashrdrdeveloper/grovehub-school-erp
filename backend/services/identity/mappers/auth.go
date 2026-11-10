package mappers

import (
	identityresponse "grovehub/services/identity/dto/response"
	"grovehub/services/identity/entity"
)

// @Mapper AuthMapper
type AuthMapper interface {
    // @Mapping AccountId=string(account.GetAccountId())
    // @Mapping PersonGid=string(account.GetPersonGid())
    // @Mapping GID=string(account.GetPersonGid())
    // @Mapping Email=account.GetEmail()
    // @Mapping Phone=account.GetPhone()
    // @Mapping Status=string(account.GetStatus())
	ToRegisterResponse(account entity.Account) identityresponse.Register
}

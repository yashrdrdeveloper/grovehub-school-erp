package mappers

import (
	identityresponse "grovehub/services/identity/dto/response"
	"grovehub/services/identity/entity"
)

// @Mapper AuthMapper
type AuthMapper interface {
    // @Mapping AccountId=string(account.AccountId)
    // @Mapping PersonGid=string(account.PersonGid)
    // @Mapping GID=string(account.PersonGid)
    // @Mapping Email=account.Email
    // @Mapping Phone=nullableString(account.Phone)
    // @Mapping Status=string(account.Status)
	ToRegisterResponse(account entity.Account) identityresponse.Register
}

package mappers

import (
	identityresponse "grovehub/services/identity/dto/response"
	"grovehub/services/identity/entity"
	identityhelpers "grovehub/services/identity/services/helpers"
)

// @Mapper AuthMapper
type AuthMapper interface {
    // @Mapping AccountId=string(account.AccountId)
    // @Mapping PersonGid=string(account.PersonGid)
    // @Mapping GID=string(account.PersonGid)
    // @Mapping Email=account.Email
    // @Mapping Phone=identityhelpers.NullableString(account.Phone)
    // @Mapping Status=string(account.Status)
	ToRegisterResponse(account entity.Account) identityresponse.Register
}

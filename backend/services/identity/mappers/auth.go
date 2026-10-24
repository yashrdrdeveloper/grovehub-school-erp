package mappers

import (
	identityresponse "grovehub/services/identity/dto/response"
	"grovehub/services/identity/entity"
)

type AuthMapper interface {
	ToRegisterResponse(account entity.Account) identityresponse.Register
}

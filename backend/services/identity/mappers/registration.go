package mappers

import (
	identityresponse "grovehub/services/identity/dto/response"
	"grovehub/services/identity/entity"
)

type RegistrationMapper interface {
	ToRegistrationResponse(account entity.Account) identityresponse.Registration
}

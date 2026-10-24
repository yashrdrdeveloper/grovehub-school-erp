package validators

import identityrequest "grovehub/services/identity/dto/request"

type AuthValidator interface {
	ValidateRegister(request identityrequest.Register) (identityrequest.Register, error)
}

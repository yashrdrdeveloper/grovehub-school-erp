package validators

import identityrequest "grovehub/services/identity/dto/request"

type RegistrationValidator interface {
	Validate(request identityrequest.Register) (identityrequest.Register, error)
}

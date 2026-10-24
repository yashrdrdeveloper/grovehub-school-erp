package impl

import (
	identityrequest "grovehub/services/identity/dto/request"
	identityvalidators "grovehub/services/identity/validators"
	validatorhelpers "grovehub/services/identity/validators/helpers"
)

type RegistrationValidator struct{}

func NewRegistrationValidator() identityvalidators.RegistrationValidator {
	return &RegistrationValidator{}
}

func (validator *RegistrationValidator) Validate(
	request identityrequest.Register,
) (identityrequest.Register, error) {
	return validatorhelpers.NormalizeRegistration(request)
}

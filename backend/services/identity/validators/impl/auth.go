package impl

import (
	identityrequest "grovehub/services/identity/dto/request"
	identityvalidators "grovehub/services/identity/validators"
	validatorhelpers "grovehub/services/identity/validators/helpers"
)

type AuthValidator struct{}

func NewAuthValidator() identityvalidators.AuthValidator {
	return &AuthValidator{}
}

func (validator *AuthValidator) ValidateRegister(
	request identityrequest.Register,
) (identityrequest.Register, error) {
	return validatorhelpers.NormalizeRegister(request)
}

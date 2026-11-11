package impl

import (
	annotationvalidation "grovehub/libs/core/annotation/validation"
	identityrequest "grovehub/services/identity/dto/request"
	identityvalidators "grovehub/services/identity/validators"
	validatorhelpers "grovehub/services/identity/validators/helpers"
)

type AuthValidator struct{}

func NewAuthValidator() identityvalidators.AuthValidator {
	return &AuthValidator{}
}

func (validator *AuthValidator) RegisterRules() annotationvalidation.Rules[identityrequest.Register] {
	return annotationvalidation.Rules[identityrequest.Register]{
		Normalize: validatorhelpers.NormalizeRegisterFields,
		Validate: validatorhelpers.ValidateRegisterCustom,
		MapDTOError: validatorhelpers.MapRegisterDTOError,
	}
}

func (validator *AuthValidator) ValidateRegister(request identityrequest.Register) (identityrequest.Register, error) {
	return annotationvalidation.Process(request, validator.RegisterRules())
}

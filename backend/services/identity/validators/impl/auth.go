package impl

import (
	annotationvalidation "grovehub/libs/core/annotation/validation"
	identityrequest "grovehub/services/identity/dto/request"
	identityvalidators "grovehub/services/identity/validators"
	validatorhelpers "grovehub/services/identity/validators/helpers"
)

type AuthValidator struct{}

func NewAuthValidator() identityvalidators.AuthValidator { return &AuthValidator{} }

func (validator *AuthValidator) RegisterRules() annotationvalidation.Rules[identityrequest.RegisterRequest] {
	return annotationvalidation.Rules[identityrequest.RegisterRequest]{
		Normalize: func(input identityrequest.RegisterRequest) (identityrequest.RegisterRequest, error) {
			body, err := validatorhelpers.NormalizeRegisterFields(input.RequestBody)
			if err != nil { return identityrequest.RegisterRequest{}, err }
			input.RequestBody = body
			return input, nil
		},
		Validate: func(input identityrequest.RegisterRequest) error {
			return validatorhelpers.ValidateRegisterCustom(input.RequestBody)
		},
		MapDTOError: validatorhelpers.MapRegisterDTOError,
	}
}

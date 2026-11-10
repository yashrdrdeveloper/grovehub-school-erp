package helpers

import (
	annotationvalidation "grovehub/libs/core/annotation/validation"

	coreerrors "grovehub/libs/core/errors"
	identityerror "grovehub/services/identity/constants/error"
	identityrequest "grovehub/services/identity/dto/request"
)

func NormalizeRegister(request identityrequest.Register) (identityrequest.Register, error) {
	if _, err := annotationvalidation.Normalize(request); err != nil {
		return identityrequest.Register{}, invalidRegistration()
	}
	if err := ValidatePassword(request.Password); err != nil {
		return identityrequest.Register{}, invalidRegistration()
	}

	email, err := NormalizeEmail(request.Email)
	if err != nil {
		return identityrequest.Register{}, invalidRegistration()
	}
	request.Email = email

	phone, err := NormalizePhone(request.PhoneNumber)
	if err != nil {
		return identityrequest.Register{}, invalidRegistration()
	}
	request.PhoneNumber = phone

	return request, nil
}

func invalidRegistration() error {
	return coreerrors.NewError(identityerror.InvalidRegistration)
}

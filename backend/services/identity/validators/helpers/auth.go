package helpers

import (
	"errors"

	annotationvalidation "grovehub/libs/core/annotation/validation"
	coreerrors "grovehub/libs/core/errors"
	identityerror "grovehub/services/identity/constants/error"
	identityrequest "grovehub/services/identity/dto/request"
)

// NormalizeRegisterFields canonicalizes technical values before DTO tag rules.
func NormalizeRegisterFields(request identityrequest.Register) (identityrequest.Register, error) {
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

// ValidateRegisterCustom checks domain rules that are not expressed by DTO tags.
func ValidateRegisterCustom(request identityrequest.Register) error {
	if err := ValidatePassword(request.Password); err != nil {
		return invalidRegistration()
	}
	return nil
}

func MapRegisterDTOError(err error) error {
	var fieldError annotationvalidation.Error
	if errors.As(err, &fieldError) {
		return invalidRegistration()
	}
	return err
}

func invalidRegistration() error {
	return coreerrors.NewError(identityerror.InvalidRegistration)
}

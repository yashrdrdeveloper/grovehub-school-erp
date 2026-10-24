package helpers

import (
	"strings"

	coreerrors "grovehub/libs/core/errors"
	identityerror "grovehub/services/identity/constants/error"
	identityrequest "grovehub/services/identity/dto/request"
)

func NormalizeRegister(request identityrequest.Register) (identityrequest.Register, error) {
	if strings.TrimSpace(request.FirstName) == "" ||
		request.MiddleName == nil ||
		strings.TrimSpace(request.LastName) == "" {
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

	if strings.TrimSpace(request.Gid) != "" {
		personGid, err := NormalizePersonGid(request.Gid)
		if err != nil {
			return identityrequest.Register{}, invalidRegistration()
		}
		request.Gid = string(personGid)
	}

	return request, nil
}

func invalidRegistration() error {
	return coreerrors.NewError(identityerror.InvalidRegistration)
}

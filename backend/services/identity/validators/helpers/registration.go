package helpers

import (
	coreerrors "grovehub/libs/core/errors"
	foundationrequest "grovehub/services/foundation/dto/request"
	identityerror "grovehub/services/identity/constants/error"
	identityrequest "grovehub/services/identity/dto/request"
)

func ValidateRegistration(request identityrequest.Register) error {
	if err := ValidatePassword(request.Password); err != nil {
		return invalidRegistration()
	}
	if countPrimaryContacts(request.Contacts.Emails) != 1 {
		return invalidRegistration()
	}
	if countPrimaryContacts(request.Contacts.Phones) > 1 {
		return invalidRegistration()
	}
	return nil
}

func countPrimaryContacts(values []foundationrequest.ContactValue) int {
	count := 0
	for _, value := range values {
		if value.IsPrimary {
			count++
		}
	}
	return count
}

func invalidRegistration() error {
	return coreerrors.NewError(identityerror.InvalidRegistration)
}

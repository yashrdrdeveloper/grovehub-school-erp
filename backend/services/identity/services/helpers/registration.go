package helpers

import (
	"errors"

	coreerrors "grovehub/libs/core/errors"
	personerror "grovehub/services/foundation/constants/error"
	foundationresponse "grovehub/services/foundation/dto/response"
	identityerror "grovehub/services/identity/constants/error"
)

var registrationDependencyErrors = map[coreerrors.Code]coreerrors.Definition{
	personerror.InvalidPersonGID.Code():          identityerror.InvalidRegistration,
	personerror.InvalidPersonRequest.Code():      identityerror.InvalidRegistration,
	personerror.InvalidPersonContact.Code():      identityerror.InvalidRegistration,
	personerror.PersonContactConflict.Code():      identityerror.RegistrationAlreadyExists,
	identityerror.InvalidAccount.Code():           identityerror.InvalidRegistration,
	identityerror.InvalidPassword.Code():          identityerror.InvalidRegistration,
	identityerror.AccountAlreadyExists.Code():     identityerror.RegistrationAlreadyExists,
	identityerror.AccountEmailAlreadyExists.Code(): identityerror.RegistrationAlreadyExists,
	identityerror.AccountPhoneAlreadyExists.Code(): identityerror.RegistrationAlreadyExists,
	identityerror.PasswordCredentialAlreadyExists.Code(): identityerror.RegistrationAlreadyExists,
}

func PrimaryRegistrationContacts(person foundationresponse.Person) (string, string, error) {
	email := ""
	for _, contact := range person.Contacts.Emails {
		if contact.IsPrimary {
			email = contact.Value
			break
		}
	}
	if email == "" {
		return "", "", coreerrors.NewError(identityerror.RegistrationFailed)
	}

	phone := ""
	for _, contact := range person.Contacts.Phones {
		if contact.IsPrimary {
			phone = contact.Value
			break
		}
	}
	return email, phone, nil
}

func MapRegistrationDependencyError(err error) error {
	var applicationError *coreerrors.Error
	if !errors.As(err, &applicationError) {
		return coreerrors.Wrap(identityerror.RegistrationFailed, err)
	}

	definition, ok := registrationDependencyErrors[applicationError.Definition().Code()]
	if !ok {
		return coreerrors.Wrap(identityerror.RegistrationFailed, err)
	}
	return coreerrors.Wrap(definition, err)
}

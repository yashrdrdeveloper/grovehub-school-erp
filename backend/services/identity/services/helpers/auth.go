package helpers

import (
	"errors"

	coreerrors "grovehub/libs/core/errors"
	coresearch "grovehub/libs/core/search"
	personerror "grovehub/services/foundation/constants/error"
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationresponse "grovehub/services/foundation/dto/response"
	identityerror "grovehub/services/identity/constants/error"
	identityrequest "grovehub/services/identity/dto/request"
)

var authDependencyErrors = map[coreerrors.Code]coreerrors.Definition{
	personerror.InvalidPersonGID.Code():                           identityerror.InvalidRegistration,
	personerror.InvalidPersonRequest.Code():                       identityerror.InvalidRegistration,
	personerror.InvalidPersonContact.Code():                       identityerror.InvalidRegistration,
	personerror.PersonNotFound.Code():                             identityerror.InvalidRegistration,
	personerror.PersonContactConflict.Code():                      identityerror.RegistrationAlreadyExists,
	identityerror.InvalidAccount.Code():                           identityerror.InvalidRegistration,
	identityerror.InvalidPassword.Code():                          identityerror.InvalidRegistration,
	identityerror.AccountAlreadyExists.Code():                     identityerror.RegistrationAlreadyExists,
	identityerror.AccountEmailAlreadyExists.Code():                identityerror.RegistrationAlreadyExists,
	identityerror.AccountPhoneAlreadyExists.Code():                identityerror.RegistrationAlreadyExists,
	identityerror.AccountRegistrationIdempotencyConflict.Code():   identityerror.RegistrationIdempotencyConflict,
	identityerror.PasswordCredentialAlreadyExists.Code():          identityerror.RegistrationAlreadyExists,
	identityerror.AccountConcurrentModification.Code():            identityerror.RegistrationFailed,
	identityerror.PasswordCredentialConcurrentModification.Code(): identityerror.RegistrationFailed,
}

func RegisterEmailSearch(email string) coresearch.Request {
	return coresearch.Request{
		Filters: []coresearch.Filter{{
			Field:    "contacts.email",
			Operator: coresearch.OperatorEquals,
			Values:   []string{email},
		}},
		PageInfo: coresearch.PageRequest{Rows: 2},
	}
}

func BuildRegisterPersonRequest(request identityrequest.Register) foundationrequest.CreatePerson {
	contacts := foundationrequest.Contacts{
		Emails: []foundationrequest.ContactValue{{Value: request.Email, IsPrimary: true}},
	}
	if request.PhoneNumber != "" {
		contacts.Phones = []foundationrequest.ContactValue{{Value: request.PhoneNumber, IsPrimary: true}}
	}
	return foundationrequest.CreatePerson{
		FirstName:  request.FirstName,
		MiddleName: *request.MiddleName,
		LastName:   request.LastName,
		Contacts:   contacts,
	}
}

func BuildRegisterPersonUpdate(request identityrequest.Register) foundationrequest.UpdatePerson {
	firstName := request.FirstName
	middleName := *request.MiddleName
	lastName := request.LastName
	return foundationrequest.UpdatePerson{
		FirstName:  &firstName,
		MiddleName: &middleName,
		LastName:   &lastName,
	}
}

func PersonHasEmail(person foundationresponse.Person, email string) bool {
	for _, contact := range person.Contacts.Emails {
		if contact.Value == email {
			return true
		}
	}
	return false
}

func MapAuthDependencyError(err error) error {
	var applicationError *coreerrors.Error
	if !errors.As(err, &applicationError) {
		return coreerrors.Wrap(identityerror.RegistrationFailed, err)
	}
	definition, ok := authDependencyErrors[applicationError.Definition().Code()]
	if !ok {
		return coreerrors.Wrap(identityerror.RegistrationFailed, err)
	}
	return coreerrors.Wrap(definition, err)
}

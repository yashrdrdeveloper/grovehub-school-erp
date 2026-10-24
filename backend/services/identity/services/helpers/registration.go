package helpers

import (
	"errors"

	searchcontract "grovehub/libs/core/contracts/search"
	coreerrors "grovehub/libs/core/errors"
	personerror "grovehub/services/foundation/constants/error"
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationresponse "grovehub/services/foundation/dto/response"
	identityerror "grovehub/services/identity/constants/error"
	identityrequest "grovehub/services/identity/dto/request"
)

var registrationDependencyErrors = map[coreerrors.Code]coreerrors.Definition{
	personerror.InvalidPersonGID.Code():                    identityerror.InvalidRegistration,
	personerror.InvalidPersonRequest.Code():                identityerror.InvalidRegistration,
	personerror.InvalidPersonContact.Code():                identityerror.InvalidRegistration,
	personerror.PersonNotFound.Code():                       identityerror.InvalidRegistration,
	personerror.PersonContactConflict.Code():                identityerror.RegistrationAlreadyExists,
	identityerror.InvalidAccount.Code():                     identityerror.InvalidRegistration,
	identityerror.InvalidPassword.Code():                    identityerror.InvalidRegistration,
	identityerror.AccountAlreadyExists.Code():               identityerror.RegistrationAlreadyExists,
	identityerror.AccountEmailAlreadyExists.Code():          identityerror.RegistrationAlreadyExists,
	identityerror.AccountPhoneAlreadyExists.Code():          identityerror.RegistrationAlreadyExists,
	identityerror.PasswordCredentialAlreadyExists.Code():    identityerror.RegistrationAlreadyExists,
}

func RegistrationEmailSearch(email string) searchcontract.Request {
	return searchcontract.Request{
		Filters: []searchcontract.Filter{
			{
				Field:    "contacts.email",
				Operator: searchcontract.OperatorEquals,
				Values:   []string{email},
			},
		},
		PageInfo: searchcontract.PageRequest{Rows: 2},
	}
}

func BuildRegistrationPersonRequest(request identityrequest.Register) foundationrequest.CreatePerson {
	contacts := foundationrequest.Contacts{
		Emails: []foundationrequest.ContactValue{
			{
				Value:     request.Email,
				IsPrimary: true,
			},
		},
	}
	if request.PhoneNumber != "" {
		contacts.Phones = []foundationrequest.ContactValue{
			{
				Value:     request.PhoneNumber,
				IsPrimary: true,
			},
		}
	}

	return foundationrequest.CreatePerson{
		FirstName:  request.FirstName,
		MiddleName: request.MiddleName,
		LastName:   request.LastName,
		Contacts:   contacts,
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

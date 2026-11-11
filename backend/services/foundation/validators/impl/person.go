package impl

import (
	"errors"
	"strings"

	annotationvalidation "grovehub/libs/core/annotation/validation"
	coreerrors "grovehub/libs/core/errors"
	coresearch "grovehub/libs/core/search"
	personerror "grovehub/services/foundation/constants/error"
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationvalidators "grovehub/services/foundation/validators"
	validatorhelpers "grovehub/services/foundation/validators/helpers"
)

type PersonValidator struct{}

func NewPersonValidator() foundationvalidators.PersonValidator {
	return &PersonValidator{}
}

func (validator *PersonValidator) CreateRules() annotationvalidation.Rules[foundationrequest.CreatePerson] {
	return annotationvalidation.Rules[foundationrequest.CreatePerson]{
		Normalize: func(request foundationrequest.CreatePerson) (foundationrequest.CreatePerson, error) {
			contacts, err := validatorhelpers.NormalizeContacts(request.Contacts)
			if err != nil {
				return foundationrequest.CreatePerson{}, err
			}
			request.Contacts = contacts
			return request, nil
		},
		Validate: func(request foundationrequest.CreatePerson) error {
			return validatorhelpers.ValidateContacts(request.Contacts)
		},
		MapDTOError: mapPersonDTOError,
	}
}

func (validator *PersonValidator) UpdateRules() annotationvalidation.Rules[foundationrequest.UpdatePerson] {
	return annotationvalidation.Rules[foundationrequest.UpdatePerson]{
		Normalize: func(request foundationrequest.UpdatePerson) (foundationrequest.UpdatePerson, error) {
			if request.Contacts.Set && request.Contacts.Value != nil {
				contacts, err := validatorhelpers.NormalizeContacts(*request.Contacts.Value)
				if err != nil {
					return foundationrequest.UpdatePerson{}, err
				}
				request.Contacts.Value = &contacts
			}
			return request, nil
		},
		Validate: func(request foundationrequest.UpdatePerson) error {
			if request.Contacts.Set && request.Contacts.Value != nil {
				return validatorhelpers.ValidateContacts(*request.Contacts.Value)
			}
			return nil
		},
		MapDTOError: mapPersonDTOError,
	}
}

func (validator *PersonValidator) ValidateCreate(request foundationrequest.CreatePerson) (foundationrequest.CreatePerson, error) {
	return annotationvalidation.Process(request, validator.CreateRules())
}

func (validator *PersonValidator) ValidateUpdate(request foundationrequest.UpdatePerson) (foundationrequest.UpdatePerson, error) {
	return annotationvalidation.Process(request, validator.UpdateRules())
}

func mapPersonDTOError(err error) error {
	var fieldError annotationvalidation.Error
	if !errors.As(err, &fieldError) {
		return err
	}
	if strings.HasPrefix(fieldError.Field, "Contacts.") || fieldError.Field == "Contacts" {
		return coreerrors.NewError(personerror.InvalidPersonContact)
	}
	return coreerrors.NewError(personerror.InvalidPersonRequest)
}

func (validator *PersonValidator) ValidateGID(gid string) error {
	return validatorhelpers.ValidateGID(gid)
}

func (validator *PersonValidator) ValidateSearch(request coresearch.Request) (coresearch.Query[string], coresearch.CursorPayload, error) {
	return validatorhelpers.NormalizePersonSearchRequest(request)
}

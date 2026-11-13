package impl

import (
	"errors"
	"strings"

	annotationvalidation "grovehub/libs/core/annotation/validation"
	coreerrors "grovehub/libs/core/errors"
	personerror "grovehub/services/foundation/constants/error"
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationvalidators "grovehub/services/foundation/validators"
	validatorhelpers "grovehub/services/foundation/validators/helpers"
)

type PersonValidator struct{}

func NewPersonValidator() foundationvalidators.PersonValidator {
	return &PersonValidator{}
}

func (validator *PersonValidator) CreateRules() annotationvalidation.Rules[foundationrequest.CreatePersonRequest] {
	return annotationvalidation.Rules[foundationrequest.CreatePersonRequest]{
		Normalize: func(request foundationrequest.CreatePersonRequest) (foundationrequest.CreatePersonRequest, error) {
			contacts, err := validatorhelpers.NormalizeContacts(request.RequestBody.Contacts)
			if err != nil { return foundationrequest.CreatePersonRequest{}, err }
			request.RequestBody.Contacts = contacts
			return request, nil
		},
		Validate: func(request foundationrequest.CreatePersonRequest) error {
			return validatorhelpers.ValidateContacts(request.RequestBody.Contacts)
		},
		MapDTOError: mapPersonDTOError,
	}
}

func (validator *PersonValidator) UpdateRules() annotationvalidation.Rules[foundationrequest.UpdatePersonRequest] {
	return annotationvalidation.Rules[foundationrequest.UpdatePersonRequest]{
		Normalize: func(request foundationrequest.UpdatePersonRequest) (foundationrequest.UpdatePersonRequest, error) {
			contacts := request.RequestBody.Contacts
			if contacts.Set && contacts.Value != nil {
				normalized, err := validatorhelpers.NormalizeContacts(*contacts.Value)
				if err != nil { return foundationrequest.UpdatePersonRequest{}, err }
				request.RequestBody.Contacts.Value = &normalized
			}
			return request, nil
		},
		Validate: func(request foundationrequest.UpdatePersonRequest) error {
			contacts := request.RequestBody.Contacts
			if contacts.Set && contacts.Value != nil {
				return validatorhelpers.ValidateContacts(*contacts.Value)
			}
			return nil
		},
		MapDTOError: mapPersonDTOError,
	}
}

func (validator *PersonValidator) SearchRules() annotationvalidation.Rules[foundationrequest.SearchPersonsRequest] {
	return annotationvalidation.Rules[foundationrequest.SearchPersonsRequest]{
		Validate: func(request foundationrequest.SearchPersonsRequest) error {
			_, _, err := validatorhelpers.NormalizePersonSearchRequest(request.RequestBody)
			return err
		},
	}
}

func mapPersonDTOError(err error) error {
	var fieldError annotationvalidation.Error
	if !errors.As(err, &fieldError) { return err }
	field := strings.TrimPrefix(fieldError.Field, "RequestBody.")
	if strings.HasPrefix(field, "Contacts.") || field == "Contacts" {
		return coreerrors.NewError(personerror.InvalidPersonContact)
	}
	return coreerrors.NewError(personerror.InvalidPersonRequest)
}

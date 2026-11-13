package impl

import (
	"errors"
	"strings"

	annotationvalidation "grovehub/libs/core/annotation/validation"
	coreerrors "grovehub/libs/core/errors"
	personerror "grovehub/services/foundation/constants/error"
	persondto "grovehub/services/foundation/dto/person"
	foundationvalidators "grovehub/services/foundation/validators"
	validatorhelpers "grovehub/services/foundation/validators/helpers"
)

type PersonValidator struct{}

func NewPersonValidator() foundationvalidators.PersonValidator {
	return &PersonValidator{}
}

func (validator *PersonValidator) CreateRules() annotationvalidation.Rules[persondto.CreatePersonRequest] {
	return annotationvalidation.Rules[persondto.CreatePersonRequest]{
		Normalize: func(request persondto.CreatePersonRequest) (persondto.CreatePersonRequest, error) {
			contacts, err := validatorhelpers.NormalizeContacts(request.RequestBody.Contacts)
			if err != nil { return persondto.CreatePersonRequest{}, err }
			request.RequestBody.Contacts = contacts
			return request, nil
		},
		Validate: func(request persondto.CreatePersonRequest) error {
			return validatorhelpers.ValidateContacts(request.RequestBody.Contacts)
		},
		MapDTOError: mapPersonDTOError,
	}
}

func (validator *PersonValidator) UpdateRules() annotationvalidation.Rules[persondto.UpdatePersonRequest] {
	return annotationvalidation.Rules[persondto.UpdatePersonRequest]{
		Normalize: func(request persondto.UpdatePersonRequest) (persondto.UpdatePersonRequest, error) {
			contacts := request.RequestBody.Contacts
			if contacts.Set && contacts.Value != nil {
				normalized, err := validatorhelpers.NormalizeContacts(*contacts.Value)
				if err != nil { return persondto.UpdatePersonRequest{}, err }
				request.RequestBody.Contacts.Value = &normalized
			}
			return request, nil
		},
		Validate: func(request persondto.UpdatePersonRequest) error {
			contacts := request.RequestBody.Contacts
			if contacts.Set && contacts.Value != nil {
				return validatorhelpers.ValidateContacts(*contacts.Value)
			}
			return nil
		},
		MapDTOError: mapPersonDTOError,
	}
}

func (validator *PersonValidator) SearchRules() annotationvalidation.Rules[persondto.SearchPersonsRequest] {
	return annotationvalidation.Rules[persondto.SearchPersonsRequest]{
		Validate: func(request persondto.SearchPersonsRequest) error {
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

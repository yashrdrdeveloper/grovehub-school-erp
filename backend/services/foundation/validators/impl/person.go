package impl

import (
	searchcontract "grovehub/libs/core/contracts/search"
	coresearch "grovehub/libs/core/search"
	foundationrequest "grovehub/services/foundation/dto/request"
	"grovehub/services/foundation/entity"
	foundationvalidators "grovehub/services/foundation/validators"
	validatorhelpers "grovehub/services/foundation/validators/helpers"
)

type PersonValidator struct{}

func NewPersonValidator() foundationvalidators.PersonValidator {
	return &PersonValidator{}
}

func (validator *PersonValidator) ValidateCreate(
	request foundationrequest.CreatePerson,
) (foundationrequest.CreatePerson, error) {
	firstName, middleName, lastName, err := validatorhelpers.NormalizePersonNames(
		request.FirstName,
		request.MiddleName,
		request.LastName,
	)
	if err != nil {
		return foundationrequest.CreatePerson{}, err
	}

	contacts, err := validatorhelpers.NormalizeContacts(request.Contacts)
	if err != nil {
		return foundationrequest.CreatePerson{}, err
	}

	request.FirstName = firstName
	request.MiddleName = middleName
	request.LastName = lastName
	request.Contacts = contacts
	return request, nil
}

func (validator *PersonValidator) ValidateUpdate(
	request foundationrequest.UpdatePerson,
) (foundationrequest.UpdatePerson, error) {
	if request.FirstName != nil {
		value, err := validatorhelpers.NormalizeRequiredName(*request.FirstName)
		if err != nil {
			return foundationrequest.UpdatePerson{}, err
		}
		request.FirstName = &value
	}
	if request.MiddleName != nil {
		value, err := validatorhelpers.NormalizeOptionalName(*request.MiddleName)
		if err != nil {
			return foundationrequest.UpdatePerson{}, err
		}
		request.MiddleName = &value
	}
	if request.LastName != nil {
		value, err := validatorhelpers.NormalizeRequiredName(*request.LastName)
		if err != nil {
			return foundationrequest.UpdatePerson{}, err
		}
		request.LastName = &value
	}
	if request.Contacts != nil {
		contacts, err := validatorhelpers.NormalizeContacts(*request.Contacts)
		if err != nil {
			return foundationrequest.UpdatePerson{}, err
		}
		request.Contacts = &contacts
	}

	return request, nil
}

func (validator *PersonValidator) ValidateGID(gid entity.GID) error {
	return validatorhelpers.ValidateGID(gid)
}

func (validator *PersonValidator) ValidateSearch(
	request searchcontract.Request,
) (coresearch.Query[entity.GID], coresearch.CursorPayload, error) {
	return validatorhelpers.NormalizePersonSearchRequest(request)
}

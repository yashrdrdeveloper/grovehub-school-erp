package impl

import (
	searchcontract "grovehub/libs/core/contracts/search"
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
) (foundationvalidators.ValidatedCreatePerson, error) {
	firstName, middleName, lastName, err := validatorhelpers.NormalizePersonNames(
		request.FirstName,
		request.MiddleName,
		request.LastName,
	)
	if err != nil {
		return foundationvalidators.ValidatedCreatePerson{}, err
	}

	contacts, err := validatorhelpers.BuildContacts(request.Contacts, nil)
	if err != nil {
		return foundationvalidators.ValidatedCreatePerson{}, err
	}

	return foundationvalidators.ValidatedCreatePerson{
		FirstName:  firstName,
		MiddleName: middleName,
		LastName:   lastName,
		Contacts:   contacts,
	}, nil
}

func (validator *PersonValidator) ValidateUpdate(
	request foundationrequest.UpdatePerson,
	existingContacts []entity.Contact,
) (foundationvalidators.ValidatedUpdatePerson, error) {
	validated := foundationvalidators.ValidatedUpdatePerson{}

	if request.FirstName != nil {
		value, err := validatorhelpers.NormalizeRequiredName(*request.FirstName)
		if err != nil {
			return foundationvalidators.ValidatedUpdatePerson{}, err
		}
		validated.FirstName = &value
	}
	if request.MiddleName != nil {
		value, err := validatorhelpers.NormalizeOptionalName(*request.MiddleName)
		if err != nil {
			return foundationvalidators.ValidatedUpdatePerson{}, err
		}
		validated.MiddleName = &value
	}
	if request.LastName != nil {
		value, err := validatorhelpers.NormalizeRequiredName(*request.LastName)
		if err != nil {
			return foundationvalidators.ValidatedUpdatePerson{}, err
		}
		validated.LastName = &value
	}
	if request.Contacts != nil {
		contacts, err := validatorhelpers.BuildContacts(*request.Contacts, existingContacts)
		if err != nil {
			return foundationvalidators.ValidatedUpdatePerson{}, err
		}
		validated.Contacts = &contacts
	}

	return validated, nil
}

func (validator *PersonValidator) ValidateGID(gid entity.GID) error {
	return validatorhelpers.ValidateGID(gid)
}

func (validator *PersonValidator) ValidateSearch(
	request searchcontract.Request,
) (foundationvalidators.ValidatedPersonSearch, error) {
	query, cursorPayload, err := validatorhelpers.NormalizePersonSearchRequest(request)
	if err != nil {
		return foundationvalidators.ValidatedPersonSearch{}, err
	}
	return foundationvalidators.ValidatedPersonSearch{
		Query:         query,
		CursorPayload: cursorPayload,
	}, nil
}

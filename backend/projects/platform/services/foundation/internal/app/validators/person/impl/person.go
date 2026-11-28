package impl

import (
	coreannotation "grovehub/libs/core/annotation"
	persondto "grovehub/services/foundation/internal/app/dto/person"
	foundationvalidators "grovehub/services/foundation/internal/app/validators/person"
	validatorhelpers "grovehub/services/foundation/internal/app/validators/person/helpers"
)

type PersonValidator struct{}

func NewPersonValidator() foundationvalidators.PersonValidator { return &PersonValidator{} }

func (validator *PersonValidator) CreateRules() coreannotation.ValidationRules[persondto.CreatePersonRequest] {
	return coreannotation.ValidationRules[persondto.CreatePersonRequest]{
		Normalize:   validatorhelpers.NormalizeCreatePersonRequest,
		Validate:    validatorhelpers.ValidateCreatePersonRequest,
		MapDTOError: validatorhelpers.MapPersonDTOError,
	}
}

func (validator *PersonValidator) UpdateRules() coreannotation.ValidationRules[persondto.UpdatePersonRequest] {
	return coreannotation.ValidationRules[persondto.UpdatePersonRequest]{
		Normalize:   validatorhelpers.NormalizeUpdatePersonRequest,
		Validate:    validatorhelpers.ValidateUpdatePersonRequest,
		MapDTOError: validatorhelpers.MapPersonDTOError,
	}
}

func (validator *PersonValidator) SearchRules() coreannotation.ValidationRules[persondto.SearchPersonsRequest] {
	return coreannotation.ValidationRules[persondto.SearchPersonsRequest]{
		Validate: validatorhelpers.ValidateSearchPersonRequest,
	}
}

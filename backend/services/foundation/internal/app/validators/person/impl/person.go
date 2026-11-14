package impl

import (
    annotationvalidation "grovehub/libs/core/annotation/validation"
    persondto "grovehub/services/foundation/client/person/dto"
    foundationvalidators "grovehub/services/foundation/internal/app/validators/person"
    validatorhelpers "grovehub/services/foundation/internal/app/validators/person/helpers"
)

type PersonValidator struct{}

func NewPersonValidator() foundationvalidators.PersonValidator { return &PersonValidator{} }

func (validator *PersonValidator) CreateRules() annotationvalidation.Rules[persondto.CreatePersonRequest] {
    return annotationvalidation.Rules[persondto.CreatePersonRequest]{
        Normalize: validatorhelpers.NormalizeCreatePersonRequest,
        Validate: validatorhelpers.ValidateCreatePersonRequest,
        MapDTOError: validatorhelpers.MapPersonDTOError,
    }
}

func (validator *PersonValidator) UpdateRules() annotationvalidation.Rules[persondto.UpdatePersonRequest] {
    return annotationvalidation.Rules[persondto.UpdatePersonRequest]{
        Normalize: validatorhelpers.NormalizeUpdatePersonRequest,
        Validate: validatorhelpers.ValidateUpdatePersonRequest,
        MapDTOError: validatorhelpers.MapPersonDTOError,
    }
}

func (validator *PersonValidator) SearchRules() annotationvalidation.Rules[persondto.SearchPersonsRequest] {
    return annotationvalidation.Rules[persondto.SearchPersonsRequest]{
        Validate: validatorhelpers.ValidateSearchPersonRequest,
    }
}

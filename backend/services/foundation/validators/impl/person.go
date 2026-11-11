package impl

import (
    "errors"
    "strings"

    annotationvalidation "grovehub/libs/core/annotation/validation"
    coreerrors "grovehub/libs/core/errors"
    coresearch "grovehub/libs/core/search"
    personerror "grovehub/services/foundation/constants/error"
    foundationrequest "grovehub/services/foundation/dto/request/person"
    foundationvalidators "grovehub/services/foundation/validators"
    validatorhelpers "grovehub/services/foundation/validators/helpers"
)

type PersonValidator struct{}

func NewPersonValidator() foundationvalidators.PersonValidator {
    return &PersonValidator{}
}

func (validator *PersonValidator) ValidateCreate(
    request foundationrequest.CreatePerson,
) (foundationrequest.CreatePerson,error) {
    normalized,err:=normalizePersonRequest(request)
    if err!=nil{return foundationrequest.CreatePerson{},err}

    contacts,err:=validatorhelpers.NormalizeContacts(normalized.Contacts)
    if err!=nil{return foundationrequest.CreatePerson{},err}

    normalized.Contacts=contacts
    return normalized,nil
}

func (validator *PersonValidator) ValidateUpdate(request foundationrequest.UpdatePerson) (foundationrequest.UpdatePerson, error) {
    normalized := request
    if normalized.FirstName.Set {
        if normalized.FirstName.Value == nil { return foundationrequest.UpdatePerson{}, coreerrors.NewError(personerror.InvalidPersonRequest) }
        name, err := validatorhelpers.NormalizeRequiredName(*normalized.FirstName.Value)
        if err != nil { return foundationrequest.UpdatePerson{}, err }
        normalized.FirstName.Value = &name
    }
    if normalized.MiddleName.Set && normalized.MiddleName.Value != nil {
        name, err := validatorhelpers.NormalizeOptionalName(*normalized.MiddleName.Value)
        if err != nil { return foundationrequest.UpdatePerson{}, err }
        normalized.MiddleName.Value = &name
    }
    if normalized.LastName.Set {
        if normalized.LastName.Value == nil { return foundationrequest.UpdatePerson{}, coreerrors.NewError(personerror.InvalidPersonRequest) }
        name, err := validatorhelpers.NormalizeRequiredName(*normalized.LastName.Value)
        if err != nil { return foundationrequest.UpdatePerson{}, err }
        normalized.LastName.Value = &name
    }
    if normalized.Contacts.Set && normalized.Contacts.Value != nil {
        contacts, err := validatorhelpers.NormalizeContacts(*normalized.Contacts.Value)
        if err != nil { return foundationrequest.UpdatePerson{}, err }
        normalized.Contacts.Value = &contacts
    }
    return normalized, nil
}

func normalizePersonRequest[T any](request T)(T,error) {
    normalized,err:=annotationvalidation.Normalize(request)
    if err==nil{return normalized,nil}

    var invalid annotationvalidation.Error
    if errors.As(err,&invalid){
        var zero T
        if strings.Contains(invalid.Field, "Contacts") {
            return zero,coreerrors.NewError(personerror.InvalidPersonContact)
        }
        return zero,coreerrors.NewError(personerror.InvalidPersonRequest)
    }
    var zero T
    return zero,err
}

func (validator *PersonValidator) ValidateGID(gid string) error {
    return validatorhelpers.ValidateGID(gid)
}

func (validator *PersonValidator) ValidateSearch(
    request coresearch.Request,
) (coresearch.Query[string],coresearch.CursorPayload,error) {
    return validatorhelpers.NormalizePersonSearchRequest(request)
}

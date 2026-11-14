package helpers

import (
    persondto "grovehub/services/foundation/client/person/dto"
)

// NormalizeCreatePersonRequest canonicalizes contact values without changing the incoming request.
func NormalizeCreatePersonRequest(request persondto.CreatePersonRequest) (persondto.CreatePersonRequest, error) {
    contacts, err := NormalizeContacts(request.RequestBody.Contacts)
    if err != nil { return persondto.CreatePersonRequest{}, err }
    request.RequestBody.Contacts = contacts
    return request, nil
}

func ValidateCreatePersonRequest(request persondto.CreatePersonRequest) error {
    return ValidateContacts(request.RequestBody.Contacts)
}

// NormalizeUpdatePersonRequest touches only explicitly supplied, non-null contacts.
func NormalizeUpdatePersonRequest(request persondto.UpdatePersonRequest) (persondto.UpdatePersonRequest, error) {
    contacts := request.RequestBody.Contacts
    if !contacts.Set || contacts.Value == nil { return request, nil }
    normalized, err := NormalizeContacts(*contacts.Value)
    if err != nil { return persondto.UpdatePersonRequest{}, err }
    request.RequestBody.Contacts.Value = &normalized
    return request, nil
}

func ValidateUpdatePersonRequest(request persondto.UpdatePersonRequest) error {
    contacts := request.RequestBody.Contacts
    if !contacts.Set || contacts.Value == nil { return nil }
    return ValidateContacts(*contacts.Value)
}

func ValidateSearchPersonRequest(request persondto.SearchPersonsRequest) error {
    _, _, err := NormalizePersonSearchRequest(request.RequestBody)
    return err
}

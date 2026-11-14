package person

import (
    corerequest "grovehub/libs/core/request"
    coresearch "grovehub/libs/core/search"
)

// The fields of each endpoint request describe only the sources that exist
// for that operation. No placeholder types or transport envelope are needed.
type CreatePersonRequest struct {
    corerequest.Request
	RequestBody CreatePerson `validate:"valid"`
}

type GetPersonByGIDRequest struct {
    corerequest.Request
	Params PersonByGID `validate:"valid"`
}

type UpdatePersonRequest struct {
    corerequest.Request
	Params PersonByGID `validate:"valid"`
	RequestBody UpdatePerson `validate:"valid"`
}

type SearchPersonsRequest struct {
    corerequest.Request
	RequestBody coresearch.Request `validate:"valid"`
}

type DeactivatePersonRequest struct {
    corerequest.Request
	Params PersonByGID `validate:"valid"`
}

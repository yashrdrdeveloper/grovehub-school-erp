package request

import coresearch "grovehub/libs/core/search"

// The fields of each endpoint request describe only the sources that exist
// for that operation. No placeholder types or transport envelope are needed.
type CreatePersonRequest struct {
	RequestBody CreatePerson `validate:"valid"`
}

type GetPersonByGIDRequest struct {
	Params PersonByGID `validate:"valid"`
}

type UpdatePersonRequest struct {
	Params PersonByGID `validate:"valid"`
	RequestBody UpdatePerson `validate:"valid"`
}

type SearchPersonsRequest struct {
	RequestBody coresearch.Request `validate:"valid"`
}

type DeactivatePersonRequest struct {
	Params PersonByGID `validate:"valid"`
}

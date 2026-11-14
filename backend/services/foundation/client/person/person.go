package person

import (
	"context"

	coresearch "grovehub/libs/core/search"
	persondto "grovehub/services/foundation/internal/app/dto/person"
)


// Client-facing aliases expose Person contracts without exporting a DTO package.
// The original types and their validation behavior stay inside internal/app/dto.
type (
    CreatePersonRequest = persondto.CreatePersonRequest
    GetPersonByGIDRequest = persondto.GetPersonByGIDRequest
    UpdatePersonRequest = persondto.UpdatePersonRequest
    SearchPersonsRequest = persondto.SearchPersonsRequest
    DeactivatePersonRequest = persondto.DeactivatePersonRequest
    CreatePerson = persondto.CreatePerson
    UpdatePerson = persondto.UpdatePerson
    PersonNames = persondto.PersonNames
    PersonByGID = persondto.PersonByGID
    ContactValue = persondto.ContactValue
    Contacts = persondto.Contacts
    ResponseContactValue = persondto.ResponseContactValue
    ResponseContacts = persondto.ResponseContacts
    Person = persondto.Person
    PersonErrorMetadata = persondto.PersonErrorMetadata
)

// @Client personClient
type Client interface {
	// @RequestLine POST /protected/persons
	CreatePerson(ctx context.Context, request persondto.CreatePersonRequest) (persondto.Person, bool, error)
	// @RequestLine GET /protected/persons/{gid}
	GetPersonByGID(ctx context.Context, request persondto.GetPersonByGIDRequest) (persondto.Person, error)
	// @RequestLine PATCH /protected/persons/{gid}
	UpdatePerson(ctx context.Context, request persondto.UpdatePersonRequest) (persondto.Person, error)
	// @RequestLine POST /protected/persons/search
	SearchPersons(ctx context.Context, request persondto.SearchPersonsRequest) (coresearch.Result[persondto.Person], error)
	// @RequestLine PATCH /protected/persons/{gid}/deactivate
	DeactivatePerson(ctx context.Context, request persondto.DeactivatePersonRequest) (persondto.Person, error)
}

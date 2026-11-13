package person

import (
	"context"

	coresearch "grovehub/libs/core/search"
	persondto "grovehub/services/foundation/dto/person"
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

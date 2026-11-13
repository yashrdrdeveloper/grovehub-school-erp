package client

import (
	"context"

	coresearch "grovehub/libs/core/search"
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationresponse "grovehub/services/foundation/dto/response"
)

// @Client personClient
type Client interface {
	// @RequestLine POST /protected/persons
	CreatePerson(ctx context.Context, request foundationrequest.CreatePersonRequest) (foundationresponse.Person, bool, error)
	// @RequestLine GET /protected/persons/{gid}
	GetPersonByGID(ctx context.Context, request foundationrequest.GetPersonByGIDRequest) (foundationresponse.Person, error)
	// @RequestLine PATCH /protected/persons/{gid}
	UpdatePerson(ctx context.Context, request foundationrequest.UpdatePersonRequest) (foundationresponse.Person, error)
	// @RequestLine POST /protected/persons/search
	SearchPersons(ctx context.Context, request foundationrequest.SearchPersonsRequest) (coresearch.Result[foundationresponse.Person], error)
	// @RequestLine PATCH /protected/persons/{gid}/deactivate
	DeactivatePerson(ctx context.Context, request foundationrequest.DeactivatePersonRequest) (foundationresponse.Person, error)
}

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
	CreatePerson(ctx context.Context, request foundationrequest.CreatePerson) (foundationresponse.Person, bool, error)
	// @RequestLine GET /protected/persons/{gid}
	GetPersonByGID(ctx context.Context, request foundationrequest.PersonByGID) (foundationresponse.Person, error)
	// @RequestLine PATCH /protected/persons/{gid}
	UpdatePerson(ctx context.Context, personRequest foundationrequest.PersonByGID, request foundationrequest.UpdatePerson) (foundationresponse.Person, error)
	// @RequestLine POST /protected/persons/search
	SearchPersons(ctx context.Context, request coresearch.Request) (coresearch.Result[foundationresponse.Person], error)
	// @RequestLine PATCH /protected/persons/{gid}/deactivate
	DeactivatePerson(ctx context.Context, request foundationrequest.PersonByGID) (foundationresponse.Person, error)
}

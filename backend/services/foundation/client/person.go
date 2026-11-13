package client

import (
	"context"

	coresearch "grovehub/libs/core/search"
	httpendpoint "grovehub/libs/core/http/endpoint"
	rest "grovehub/libs/core/pipeline/rest"
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationresponse "grovehub/services/foundation/dto/response"
)

// @Client personClient
type Client interface {
	// @RequestLine POST /protected/persons
	CreatePerson(ctx context.Context, request httpendpoint.Request[rest.Empty, rest.Empty, foundationrequest.CreatePerson]) (foundationresponse.Person, bool, error)
	// @RequestLine GET /protected/persons/{gid}
	GetPersonByGID(ctx context.Context, request httpendpoint.Request[foundationrequest.PersonByGID, rest.Empty, rest.Empty]) (foundationresponse.Person, error)
	// @RequestLine PATCH /protected/persons/{gid}
	UpdatePerson(ctx context.Context, request httpendpoint.Request[foundationrequest.PersonByGID, rest.Empty, foundationrequest.UpdatePerson]) (foundationresponse.Person, error)
	// @RequestLine POST /protected/persons/search
	SearchPersons(ctx context.Context, request httpendpoint.Request[rest.Empty, rest.Empty, coresearch.Request]) (coresearch.Result[foundationresponse.Person], error)
	// @RequestLine PATCH /protected/persons/{gid}/deactivate
	DeactivatePerson(ctx context.Context, request httpendpoint.Request[foundationrequest.PersonByGID, rest.Empty, rest.Empty]) (foundationresponse.Person, error)
}

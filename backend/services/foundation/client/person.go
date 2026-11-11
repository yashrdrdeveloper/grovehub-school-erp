package client

import (
	"context"

	coresearch "grovehub/libs/core/search"
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationresponse "grovehub/services/foundation/dto/response"
)

// @Client personClient
type Client interface {
	// @Proxy personService.Create
	CreatePerson(
		ctx context.Context,
		request foundationrequest.CreatePerson,
	) (foundationresponse.Person, bool, error)
	// @Proxy personService.GetByGID
	GetPersonByGID(ctx context.Context, request foundationrequest.PersonByGID) (foundationresponse.Person, error)
	// @Proxy personService.Update
	UpdatePerson(
		ctx context.Context,
		personRequest foundationrequest.PersonByGID,
		request foundationrequest.UpdatePerson,
	) (foundationresponse.Person, error)
	// @Proxy personService.Search
	SearchPersons(
		ctx context.Context,
		request coresearch.Request,
	) (coresearch.Result[foundationresponse.Person], error)
	// @Proxy personService.Deactivate
	DeactivatePerson(ctx context.Context, request foundationrequest.PersonByGID) (foundationresponse.Person, error)
}

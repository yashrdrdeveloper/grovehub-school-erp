package client

import (
	"context"

	coresearch "grovehub/libs/core/search"
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationresponse "grovehub/services/foundation/dto/response"
)

type Client interface {
	CreatePerson(
		ctx context.Context,
		request foundationrequest.CreatePerson,
	) (foundationresponse.Person, bool, error)
	GetPersonByGID(ctx context.Context, request foundationrequest.PersonByGID) (foundationresponse.Person, error)
	UpdatePerson(
		ctx context.Context,
		personRequest foundationrequest.PersonByGID,
		request foundationrequest.UpdatePerson,
	) (foundationresponse.Person, error)
	SearchPersons(
		ctx context.Context,
		request coresearch.Request,
	) (coresearch.Result[foundationresponse.Person], error)
	DeactivatePerson(ctx context.Context, request foundationrequest.PersonByGID) (foundationresponse.Person, error)
}

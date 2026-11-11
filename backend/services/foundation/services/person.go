package services

import (
	"context"

	coresearch "grovehub/libs/core/search"
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationresponse "grovehub/services/foundation/dto/response"
)

type PersonService interface {
	Create(
		ctx context.Context,
		request foundationrequest.CreatePerson,
	) (foundationresponse.Person, bool, error)
	GetByGID(ctx context.Context, request foundationrequest.PersonByGID) (foundationresponse.Person, error)
	Update(
		ctx context.Context,
		personRequest foundationrequest.PersonByGID,
		request foundationrequest.UpdatePerson,
	) (foundationresponse.Person, error)
	Search(
		ctx context.Context,
		request coresearch.Request,
	) (coresearch.Result[foundationresponse.Person], error)
	Deactivate(ctx context.Context, request foundationrequest.PersonByGID) (foundationresponse.Person, error)
}

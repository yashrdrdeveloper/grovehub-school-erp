package services

import (
	"context"

	coresearch "grovehub/libs/core/search"
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationresponse "grovehub/services/foundation/dto/response"
)

type PersonService interface {
	Create(ctx context.Context, request foundationrequest.CreatePersonRequest) (foundationresponse.Person, bool, error)
	GetByGID(ctx context.Context, request foundationrequest.GetPersonByGIDRequest) (foundationresponse.Person, error)
	Update(ctx context.Context, request foundationrequest.UpdatePersonRequest) (foundationresponse.Person, error)
	Search(ctx context.Context, request foundationrequest.SearchPersonsRequest) (coresearch.Result[foundationresponse.Person], error)
	Deactivate(ctx context.Context, request foundationrequest.DeactivatePersonRequest) (foundationresponse.Person, error)
}

package services

import (
	"context"

	coresearch "grovehub/libs/core/search"
	persondto "grovehub/services/foundation/dto/person"
)

type PersonService interface {
	Create(ctx context.Context, request persondto.CreatePersonRequest) (persondto.Person, bool, error)
	GetByGID(ctx context.Context, request persondto.GetPersonByGIDRequest) (persondto.Person, error)
	Update(ctx context.Context, request persondto.UpdatePersonRequest) (persondto.Person, error)
	Search(ctx context.Context, request persondto.SearchPersonsRequest) (coresearch.Result[persondto.Person], error)
	Deactivate(ctx context.Context, request persondto.DeactivatePersonRequest) (persondto.Person, error)
}

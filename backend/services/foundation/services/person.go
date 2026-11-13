package services

import (
	"context"

	corerequest "grovehub/libs/core/request"
	coresearch "grovehub/libs/core/search"
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationresponse "grovehub/services/foundation/dto/response"
)

type PersonService interface {
	Create(ctx context.Context, request corerequest.Request[corerequest.Empty, corerequest.Empty, foundationrequest.CreatePerson]) (foundationresponse.Person, bool, error)
	GetByGID(ctx context.Context, request corerequest.Request[foundationrequest.PersonByGID, corerequest.Empty, corerequest.Empty]) (foundationresponse.Person, error)
	Update(ctx context.Context, request corerequest.Request[foundationrequest.PersonByGID, corerequest.Empty, foundationrequest.UpdatePerson]) (foundationresponse.Person, error)
	Search(ctx context.Context, request corerequest.Request[corerequest.Empty, corerequest.Empty, coresearch.Request]) (coresearch.Result[foundationresponse.Person], error)
	Deactivate(ctx context.Context, request corerequest.Request[foundationrequest.PersonByGID, corerequest.Empty, corerequest.Empty]) (foundationresponse.Person, error)
}

package controller

import (
	"context"
	"net/http"

	searchcontract "grovehub/libs/core/contracts/search"
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
		request searchcontract.Request,
	) (searchcontract.Result[foundationresponse.Person], error)
	Deactivate(ctx context.Context, request foundationrequest.PersonByGID) (foundationresponse.Person, error)
}

type PersonController interface {
	Create(writer http.ResponseWriter, request *http.Request)
	GetByGID(writer http.ResponseWriter, request *http.Request)
	Update(writer http.ResponseWriter, request *http.Request)
	Search(writer http.ResponseWriter, request *http.Request)
	Deactivate(writer http.ResponseWriter, request *http.Request)
}

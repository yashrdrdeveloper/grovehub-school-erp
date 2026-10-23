package controller

import (
	"context"
	"net/http"

	httpresponse "grovehub/libs/core/http/response"
	personerror "grovehub/services/foundation/constants/error"
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationresponse "grovehub/services/foundation/dto/response"
)

type PersonService interface {
	Create(ctx context.Context) (foundationresponse.Person, error)
	GetByGID(ctx context.Context, request foundationrequest.PersonByGID) (foundationresponse.Person, error)
	Deactivate(ctx context.Context, request foundationrequest.PersonByGID) (foundationresponse.Person, error)
}

type PersonController struct {
	service PersonService
}

func NewPersonController(service PersonService) *PersonController {
	return &PersonController{service: service}
}

func (controller *PersonController) Create(writer http.ResponseWriter, request *http.Request) {
	person, err := controller.service.Create(request.Context())
	if err != nil {
		writePersonError(
			request.Context(),
			writer,
			err,
			personerror.PersonCreateFailed,
		)
		return
	}

	_ = httpresponse.WriteSuccess(writer, http.StatusCreated, person)
}

func (controller *PersonController) GetByGID(writer http.ResponseWriter, request *http.Request) {
	personRequest := foundationrequest.PersonByGID{
		GID: request.PathValue("gid"),
	}

	person, err := controller.service.GetByGID(request.Context(), personRequest)
	if err != nil {
		writePersonError(
			request.Context(),
			writer,
			err,
			personerror.PersonReadFailed,
		)
		return
	}

	_ = httpresponse.WriteSuccess(writer, http.StatusOK, person)
}

func (controller *PersonController) Deactivate(writer http.ResponseWriter, request *http.Request) {
	personRequest := foundationrequest.PersonByGID{
		GID: request.PathValue("gid"),
	}

	person, err := controller.service.Deactivate(request.Context(), personRequest)
	if err != nil {
		writePersonError(
			request.Context(),
			writer,
			err,
			personerror.PersonUpdateFailed,
		)
		return
	}

	_ = httpresponse.WriteSuccess(writer, http.StatusOK, person)
}

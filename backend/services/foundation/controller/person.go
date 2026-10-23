package controller

import (
	"context"
	"errors"
	"net/http"

	coreerrors "grovehub/libs/core/errors"
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
		controller.writeError(
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
		controller.writeError(
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
		controller.writeError(
			request.Context(),
			writer,
			err,
			personerror.PersonUpdateFailed,
		)
		return
	}

	_ = httpresponse.WriteSuccess(writer, http.StatusOK, person)
}

func (controller *PersonController) writeError(
	ctx context.Context,
	writer http.ResponseWriter,
	err error,
	fallback coreerrors.Definition,
) {
	var applicationError *coreerrors.Error
	if !errors.As(err, &applicationError) {
		_ = httpresponse.WriteErrorContext(
			ctx,
			writer,
			http.StatusInternalServerError,
			fallback,
		)
		return
	}

	statusCode := controller.errorHTTPStatus(applicationError.Definition())
	metadata := applicationError.Metadata()
	if metadata == nil {
		_ = httpresponse.WriteErrorContext(
			ctx,
			writer,
			statusCode,
			applicationError.Definition(),
		)
		return
	}

	_ = httpresponse.WriteErrorWithMetadataContext(
		ctx,
		writer,
		statusCode,
		applicationError.Definition(),
		metadata,
	)
}

func (controller *PersonController) errorHTTPStatus(definition coreerrors.Definition) int {
	switch definition.Code() {
	case personerror.InvalidPersonGID.Code():
		return http.StatusBadRequest
	case personerror.PersonNotFound.Code():
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

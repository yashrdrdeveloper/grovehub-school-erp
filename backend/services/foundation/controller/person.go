package controller

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	searchcontract "grovehub/libs/core/contracts/search"
	coreerrors "grovehub/libs/core/errors"
	httpresponse "grovehub/libs/core/http/response"
	personerror "grovehub/services/foundation/constants/error"
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

type PersonController struct {
	service PersonService
}

func NewPersonController(service PersonService) *PersonController {
	return &PersonController{service: service}
}

func (controller *PersonController) Create(writer http.ResponseWriter, request *http.Request) {
	var personRequest foundationrequest.CreatePerson
	if err := controller.decodeRequest(request, &personRequest); err != nil {
		_ = httpresponse.WriteErrorContext(request.Context(), writer, http.StatusBadRequest, personerror.InvalidPersonRequest)
		return
	}

	person, created, err := controller.service.Create(request.Context(), personRequest)
	if err != nil {
		controller.writeError(request.Context(), writer, err, personerror.PersonCreateFailed)
		return
	}

	statusCode := http.StatusOK
	if created {
		statusCode = http.StatusCreated
	}
	_ = httpresponse.WriteSuccess(writer, statusCode, person)
}

func (controller *PersonController) GetByGID(writer http.ResponseWriter, request *http.Request) {
	personRequest := foundationrequest.PersonByGID{GID: request.PathValue("gid")}
	person, err := controller.service.GetByGID(request.Context(), personRequest)
	if err != nil {
		controller.writeError(request.Context(), writer, err, personerror.PersonReadFailed)
		return
	}
	_ = httpresponse.WriteSuccess(writer, http.StatusOK, person)
}

func (controller *PersonController) Update(writer http.ResponseWriter, request *http.Request) {
	personRequest := foundationrequest.PersonByGID{GID: request.PathValue("gid")}
	var updateRequest foundationrequest.UpdatePerson
	if err := controller.decodeRequest(request, &updateRequest); err != nil {
		_ = httpresponse.WriteErrorContext(request.Context(), writer, http.StatusBadRequest, personerror.InvalidPersonRequest)
		return
	}
	person, err := controller.service.Update(request.Context(), personRequest, updateRequest)
	if err != nil {
		controller.writeError(request.Context(), writer, err, personerror.PersonUpdateFailed)
		return
	}
	_ = httpresponse.WriteSuccess(writer, http.StatusOK, person)
}

func (controller *PersonController) Search(writer http.ResponseWriter, request *http.Request) {
	var searchRequest searchcontract.Request
	if err := controller.decodeRequest(request, &searchRequest); err != nil {
		_ = httpresponse.WriteErrorContext(request.Context(), writer, http.StatusBadRequest, personerror.InvalidPersonSearch)
		return
	}
	result, err := controller.service.Search(request.Context(), searchRequest)
	if err != nil {
		controller.writeError(request.Context(), writer, err, personerror.PersonSearchFailed)
		return
	}
	_ = httpresponse.WriteSuccess(writer, http.StatusOK, result)
}

func (controller *PersonController) Deactivate(writer http.ResponseWriter, request *http.Request) {
	personRequest := foundationrequest.PersonByGID{GID: request.PathValue("gid")}
	person, err := controller.service.Deactivate(request.Context(), personRequest)
	if err != nil {
		controller.writeError(request.Context(), writer, err, personerror.PersonUpdateFailed)
		return
	}
	_ = httpresponse.WriteSuccess(writer, http.StatusOK, person)
}

func (controller *PersonController) decodeRequest(request *http.Request, destination any) error {
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain one JSON value")
	}
	return nil
}

func (controller *PersonController) writeError(ctx context.Context, writer http.ResponseWriter, err error, fallback coreerrors.Definition) {
	var applicationError *coreerrors.Error
	if !errors.As(err, &applicationError) {
		_ = httpresponse.WriteErrorContext(ctx, writer, http.StatusInternalServerError, fallback)
		return
	}
	statusCode := controller.errorHTTPStatus(applicationError.Definition())
	metadata := applicationError.Metadata()
	if metadata == nil {
		_ = httpresponse.WriteErrorContext(ctx, writer, statusCode, applicationError.Definition())
		return
	}
	_ = httpresponse.WriteErrorWithMetadataContext(ctx, writer, statusCode, applicationError.Definition(), metadata)
}

func (controller *PersonController) errorHTTPStatus(definition coreerrors.Definition) int {
	switch definition.Code() {
	case personerror.InvalidPersonGID.Code(), personerror.InvalidPersonRequest.Code(), personerror.InvalidPersonContact.Code(), personerror.InvalidPersonSearch.Code(), personerror.InvalidPersonSearchCursor.Code():
		return http.StatusBadRequest
	case personerror.PersonNotFound.Code():
		return http.StatusNotFound
	case personerror.PersonContactConflict.Code():
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

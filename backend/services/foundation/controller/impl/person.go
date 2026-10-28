package impl

import (
	"net/http"

	httpresponse "grovehub/libs/core/http/response"
	coresearch "grovehub/libs/core/search"
	personerror "grovehub/services/foundation/constants/error"
	"grovehub/services/foundation/controller"
	controllerhelpers "grovehub/services/foundation/controller/helpers"
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationservices "grovehub/services/foundation/services"
)

type PersonController struct {
	service foundationservices.PersonService
}

func NewPersonController(service foundationservices.PersonService) controller.PersonController {
	return &PersonController{service: service}
}

func (controller *PersonController) Create(writer http.ResponseWriter, request *http.Request) {
	var personRequest foundationrequest.CreatePerson
	if err := controllerhelpers.DecodeRequest(request, &personRequest); err != nil {
		_ = httpresponse.WriteErrorContext(request.Context(), writer, http.StatusBadRequest, personerror.InvalidPersonRequest)
		return
	}

	person, created, err := controller.service.Create(request.Context(), personRequest)
	if err != nil {
		controllerhelpers.WritePersonError(request.Context(), writer, err, personerror.PersonCreateFailed)
		return
	}

	statusCode := http.StatusOK
	if created {
		statusCode = http.StatusCreated
	}
	_ = httpresponse.WriteSuccess(writer, statusCode, person)
}

func (controller *PersonController) GetByGID(writer http.ResponseWriter, request *http.Request) {
	person, err := controller.service.GetByGID(
		request.Context(),
		foundationrequest.PersonByGID{GID: request.PathValue("gid")},
	)
	if err != nil {
		controllerhelpers.WritePersonError(request.Context(), writer, err, personerror.PersonReadFailed)
		return
	}

	_ = httpresponse.WriteSuccess(writer, http.StatusOK, person)
}

func (controller *PersonController) Update(writer http.ResponseWriter, request *http.Request) {
	var updateRequest foundationrequest.UpdatePerson
	if err := controllerhelpers.DecodeRequest(request, &updateRequest); err != nil {
		_ = httpresponse.WriteErrorContext(request.Context(), writer, http.StatusBadRequest, personerror.InvalidPersonRequest)
		return
	}

	person, err := controller.service.Update(
		request.Context(),
		foundationrequest.PersonByGID{GID: request.PathValue("gid")},
		updateRequest,
	)
	if err != nil {
		controllerhelpers.WritePersonError(request.Context(), writer, err, personerror.PersonUpdateFailed)
		return
	}

	_ = httpresponse.WriteSuccess(writer, http.StatusOK, person)
}

func (controller *PersonController) Search(writer http.ResponseWriter, request *http.Request) {
	var searchRequest coresearch.Request
	if err := controllerhelpers.DecodeRequest(request, &searchRequest); err != nil {
		_ = httpresponse.WriteErrorContext(request.Context(), writer, http.StatusBadRequest, personerror.InvalidPersonSearch)
		return
	}

	result, err := controller.service.Search(request.Context(), searchRequest)
	if err != nil {
		controllerhelpers.WritePersonError(request.Context(), writer, err, personerror.PersonSearchFailed)
		return
	}

	_ = httpresponse.WriteSuccess(writer, http.StatusOK, result)
}

func (controller *PersonController) Deactivate(writer http.ResponseWriter, request *http.Request) {
	person, err := controller.service.Deactivate(
		request.Context(),
		foundationrequest.PersonByGID{GID: request.PathValue("gid")},
	)
	if err != nil {
		controllerhelpers.WritePersonError(request.Context(), writer, err, personerror.PersonUpdateFailed)
		return
	}

	_ = httpresponse.WriteSuccess(writer, http.StatusOK, person)
}

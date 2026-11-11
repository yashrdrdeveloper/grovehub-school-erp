package impl

import (
	"net/http"

	httpresponse "grovehub/libs/core/http/response"
	rest "grovehub/libs/core/pipeline/rest"
	annotationvalidation "grovehub/libs/core/annotation/validation"
	coresearch "grovehub/libs/core/search"
	personerror "grovehub/services/foundation/constants/error"
	protected "grovehub/services/foundation/controller/protected"
	controllerhelpers "grovehub/services/foundation/controller/protected/helpers"
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationservices "grovehub/services/foundation/services"
	foundationvalidators "grovehub/services/foundation/validators"
)

type PersonController struct {
	service   foundationservices.PersonService
	validator foundationvalidators.PersonValidator
}

func NewPersonController(service foundationservices.PersonService, validator foundationvalidators.PersonValidator) protected.PersonController {
	return &PersonController{service: service, validator: validator}
}

func (controller *PersonController) Create(writer http.ResponseWriter, request *http.Request) {
	personRequest, err := rest.Run(request, rest.Config[foundationrequest.CreatePerson]{
		Rules: controller.validator.CreateRules(),
	})
	if err != nil {
		controllerhelpers.WritePersonRequestError(request.Context(), writer, err, personerror.InvalidPersonRequest)
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
	person, err := controller.service.GetByGID(request.Context(), foundationrequest.PersonByGID{GID: request.PathValue("gid")})
	if err != nil {
		controllerhelpers.WritePersonError(request.Context(), writer, err, personerror.PersonReadFailed)
		return
	}
	_ = httpresponse.WriteSuccess(writer, http.StatusOK, person)
}

func (controller *PersonController) Update(writer http.ResponseWriter, request *http.Request) {
	updateRequest, err := rest.Run(request, rest.Config[foundationrequest.UpdatePerson]{
		Rules: controller.validator.UpdateRules(),
	})
	if err != nil {
		controllerhelpers.WritePersonRequestError(request.Context(), writer, err, personerror.InvalidPersonRequest)
		return
	}
	person, err := controller.service.Update(
		request.Context(), foundationrequest.PersonByGID{GID: request.PathValue("gid")}, updateRequest,
	)
	if err != nil {
		controllerhelpers.WritePersonError(request.Context(), writer, err, personerror.PersonUpdateFailed)
		return
	}
	_ = httpresponse.WriteSuccess(writer, http.StatusOK, person)
}

func (controller *PersonController) Search(writer http.ResponseWriter, request *http.Request) {
	searchRequest, err := rest.Run(request, rest.Config[coresearch.Request]{
		Rules: annotationvalidation.Rules[coresearch.Request]{
			Validate: func(value coresearch.Request) error {
				_, _, err := controller.validator.ValidateSearch(value)
				return err
			},
		},
	})
	if err != nil {
		controllerhelpers.WritePersonRequestError(request.Context(), writer, err, personerror.InvalidPersonSearch)
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
	person, err := controller.service.Deactivate(request.Context(), foundationrequest.PersonByGID{GID: request.PathValue("gid")})
	if err != nil {
		controllerhelpers.WritePersonError(request.Context(), writer, err, personerror.PersonUpdateFailed)
		return
	}
	_ = httpresponse.WriteSuccess(writer, http.StatusOK, person)
}

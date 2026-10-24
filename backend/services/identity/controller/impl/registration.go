package impl

import (
	"net/http"

	httprequest "grovehub/libs/core/http/request"
	httpresponse "grovehub/libs/core/http/response"
	identityerror "grovehub/services/identity/constants/error"
	"grovehub/services/identity/controller"
	controllerhelpers "grovehub/services/identity/controller/helpers"
	identityrequest "grovehub/services/identity/dto/request"
	identityservices "grovehub/services/identity/services"
)

type RegistrationController struct {
	service identityservices.RegistrationService
}

func NewRegistrationController(service identityservices.RegistrationService) controller.RegistrationController {
	return &RegistrationController{service: service}
}

func (controller *RegistrationController) Register(writer http.ResponseWriter, request *http.Request) {
	var registrationRequest identityrequest.Register
	if err := httprequest.DecodeJSON(request, &registrationRequest); err != nil {
		_ = httpresponse.WriteErrorContext(
			request.Context(),
			writer,
			http.StatusBadRequest,
			identityerror.InvalidRegistration,
		)
		return
	}

	registration, err := controller.service.Register(request.Context(), registrationRequest)
	if err != nil {
		controllerhelpers.WriteRegistrationError(request.Context(), writer, err)
		return
	}

	_ = httpresponse.WriteSuccess(writer, http.StatusCreated, registration)
}

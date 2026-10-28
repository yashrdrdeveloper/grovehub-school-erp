package impl

import (
	"net/http"

	httpheaders "grovehub/libs/core/http/headers"
	httprequest "grovehub/libs/core/http/request"
	httpresponse "grovehub/libs/core/http/response"
	identityerror "grovehub/services/identity/constants/error"
	"grovehub/services/identity/controller"
	controllerhelpers "grovehub/services/identity/controller/helpers"
	identityrequest "grovehub/services/identity/dto/request"
	identityservices "grovehub/services/identity/services"
)

type AuthController struct {
	service identityservices.AuthService
}

func NewAuthController(service identityservices.AuthService) controller.AuthController {
	return &AuthController{service: service}
}

func (controller *AuthController) Register(writer http.ResponseWriter, request *http.Request) {
	var registerRequest identityrequest.Register
	if err := httprequest.DecodeJSON(request, &registerRequest); err != nil {
		_ = httpresponse.WriteErrorContext(
			request.Context(),
			writer,
			http.StatusBadRequest,
			identityerror.InvalidRegistration,
		)
		return
	}
	registerRequest.IdempotencyKey = request.Header.Get(httpheaders.IdempotencyKey)

	registration, err := controller.service.Register(request.Context(), registerRequest)
	if err != nil {
		controllerhelpers.WriteRegisterError(request.Context(), writer, err)
		return
	}

	_ = httpresponse.WriteSuccess(writer, http.StatusCreated, registration)
}

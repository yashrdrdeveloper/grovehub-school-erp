package impl

import (
	"errors"
	"net/http"

	httprequest "grovehub/libs/core/http/request"
	httpresponse "grovehub/libs/core/http/response"
	rest "grovehub/libs/core/pipeline/rest"
	identityerror "grovehub/services/identity/constants/error"
	"grovehub/services/identity/controller"
	controllerhelpers "grovehub/services/identity/controller/helpers"
	identityrequest "grovehub/services/identity/dto/request"
	identityservices "grovehub/services/identity/services"
	identityvalidators "grovehub/services/identity/validators"
)

type AuthController struct {
	service   identityservices.AuthService
	validator identityvalidators.AuthValidator
}

func NewAuthController(service identityservices.AuthService, validator identityvalidators.AuthValidator) controller.AuthController {
	return &AuthController{service: service, validator: validator}
}

func (controller *AuthController) Register(writer http.ResponseWriter, request *http.Request) {
	registerRequest, err := rest.Run(request, rest.Config[identityrequest.Register]{
		Decode: httprequest.DecodeStrictJSON,
		Normalize: controller.validator.ValidateRegister,
	})
	if err != nil {
		var decodeError *rest.DecodeError
		if errors.As(err, &decodeError) {
			_ = httpresponse.WriteErrorContext(request.Context(), writer, http.StatusBadRequest, identityerror.InvalidRegistration)
		} else {
			controllerhelpers.WriteRegisterError(request.Context(), writer, err)
		}
		return
	}

	registration, err := controller.service.Register(request.Context(), registerRequest)
	if err != nil {
		controllerhelpers.WriteRegisterError(request.Context(), writer, err)
		return
	}

	_ = httpresponse.WriteSuccess(writer, http.StatusCreated, registration)
}

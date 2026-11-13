package impl

import (
	"context"
	"net/http"

	httprequest "grovehub/libs/core/http/request"
	rest "grovehub/libs/core/pipeline/rest"
	identityerror "grovehub/services/identity/constants/error"
	"grovehub/services/identity/controller"
	identityrequest "grovehub/services/identity/dto/request"
	identityresponse "grovehub/services/identity/dto/response"
	identityservices "grovehub/services/identity/services"
	identityvalidators "grovehub/services/identity/validators"
)

type AuthController struct {
	service identityservices.AuthService
	validator identityvalidators.AuthValidator
}

func NewAuthController(service identityservices.AuthService, validator identityvalidators.AuthValidator) controller.AuthController {
	return &AuthController{service: service, validator: validator}
}

func (controller *AuthController) Register(writer http.ResponseWriter, request *http.Request) {
	rest.Run(writer, request, rest.Config[rest.Empty, rest.Empty, identityrequest.Register, identityresponse.Register]{
		Decode: httprequest.DecodeStrictJSON,
		Body: controller.validator.RegisterRules(),
		InvalidBody: identityerror.InvalidRegistration,
		Fallback: identityerror.RegistrationFailed,
		Service: func(ctx context.Context, input rest.Input[rest.Empty, rest.Empty, identityrequest.Register]) (rest.Result[identityresponse.Register], error) {
			registration, err := controller.service.Register(ctx, input.Body)
			return rest.Created(registration), err
		},
	})
}

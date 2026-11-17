package impl

import (
    "context"
    "net/http"
    "strings"

    httprequest "grovehub/libs/core/http/request"
    rest "grovehub/libs/core/pipeline/rest"
    identityerror "grovehub/services/identity/constants/error"
    protected "grovehub/services/identity/controller/protected"
    identityrequest "grovehub/services/identity/dto/request"
    identityresponse "grovehub/services/identity/dto/response"
    identityservices "grovehub/services/identity/services"
    identityvalidators "grovehub/services/identity/validators"
)

type AuthController struct {
    service identityservices.AuthService
    validator identityvalidators.AuthValidator
}

func NewAuthController(service identityservices.AuthService, validator identityvalidators.AuthValidator) protected.AuthController {
    return &AuthController{service: service, validator: validator}
}

func (controller *AuthController) Register(writer http.ResponseWriter, request *http.Request) {
    rest.RunStruct(writer, request, rest.StructConfig[identityrequest.RegisterRequest, identityresponse.Register]{
        Decode: httprequest.DecodeStrictJSON,
        Rules: controller.validator.RegisterRules(),
        InvalidRequestBody: identityerror.InvalidRegistration,
        Fallback: identityerror.RegistrationFailed,
        Service: func(ctx context.Context, input identityrequest.RegisterRequest) (rest.Result[identityresponse.Register], error) {
            registration, err := controller.service.Register(ctx, input)
            return rest.Created(registration), err
        },
    })
}

func (controller *AuthController) Logout(writer http.ResponseWriter, request *http.Request) {
    accessToken := strings.TrimPrefix(request.Header.Get("Authorization"), "Bearer ")
    rest.RunStruct(writer, request, rest.StructConfig[identityrequest.LogoutRequest, identityresponse.Logout]{
        Fallback: identityerror.AuthenticationUnavailable,
        Service: func(ctx context.Context, input identityrequest.LogoutRequest) (rest.Result[identityresponse.Logout], error) {
            err := controller.service.Logout(ctx, accessToken)
            return rest.OK(identityresponse.Logout{Status: "LOGGED_OUT"}), err
        },
    })
}

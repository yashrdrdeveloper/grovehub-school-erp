package impl

import (
    "context"
    "net/http"

    httprequest "grovehub/libs/core/http/request"
    rest "grovehub/libs/core/pipeline/rest"
    identityerror "grovehub/services/identity/constants/error"
    public "grovehub/services/identity/controller/public"
    identityrequest "grovehub/services/identity/dto/request"
    identityresponse "grovehub/services/identity/dto/response"
    identityservices "grovehub/services/identity/services"
)

type AuthController struct { service identityservices.AuthService }

func NewAuthController(service identityservices.AuthService) public.AuthController { return &AuthController{service: service} }

func (controller *AuthController) Login(writer http.ResponseWriter, request *http.Request) {
    rest.RunStruct(writer, request, rest.StructConfig[identityrequest.LoginRequest, identityresponse.Login]{
        Decode: httprequest.DecodeStrictJSON,
        InvalidRequestBody: identityerror.InvalidLogin,
        Fallback: identityerror.AuthenticationUnavailable,
        Service: func(ctx context.Context, input identityrequest.LoginRequest) (rest.Result[identityresponse.Login], error) {
            result, err := controller.service.Login(ctx, input)
            return rest.Result[identityresponse.Login]{Data: result}, err
        },
    })
}

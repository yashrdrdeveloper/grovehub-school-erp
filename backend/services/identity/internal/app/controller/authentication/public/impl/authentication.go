package impl

import (
	"context"
	"net/http"

	httprequest "grovehub/libs/core/http/request"
	rest "grovehub/libs/core/pipeline/rest"
	definitions "grovehub/services/identity/internal/app/constants/error"
	"grovehub/services/identity/internal/app/controller/helpers"
	public "grovehub/services/identity/internal/app/controller/authentication/public"
	dto "grovehub/services/identity/internal/app/dto/authentication"
	auth "grovehub/services/identity/internal/app/services/authentication"
	secrets "grovehub/services/identity/internal/app/services/authentication/helpers"
	"grovehub/services/identity/internal/app/transport/cookies"
)

type Controller struct {
	cookies cookies.Manager
	service auth.Service
}

var _ public.AuthenticationController = (*Controller)(nil)

func New(service auth.Service, manager cookies.Manager) public.AuthenticationController {
	return &Controller{service: service, cookies: manager}
}

func (controller *Controller) Login(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[dto.LoginRequest, dto.LoginResult]{
		Decode: httprequest.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.LoginRequest) (rest.Result[dto.LoginResult], error) {
			result, err := controller.service.Login(ctx, input.RequestBody, request.UserAgent())
			if err != nil {
				return rest.Result[dto.LoginResult]{}, helpers.Map(err)
			}
			csrf, err := secrets.NewSecret(32)
			if err != nil {
				return rest.Result[dto.LoginResult]{}, helpers.Map(err)
			}
			controller.cookies.Issue(writer, result.Cookies.Access, result.Cookies.Refresh, csrf)
			return rest.OK(result.Data), nil
		},
	})
}

func (controller *Controller) Refresh(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[dto.EmptyRequest, dto.RefreshResult]{
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, _ dto.EmptyRequest) (rest.Result[dto.RefreshResult], error) {
			result, err := controller.service.Refresh(ctx, controller.cookies.Read(request, cookies.RefreshCookie))
			if err != nil {
				return rest.Result[dto.RefreshResult]{}, helpers.Map(err)
			}
			controller.cookies.Issue(writer, result.Cookies.Access, result.Cookies.Refresh, "")
			return rest.OK(result.Data), nil
		},
	})
}

func (controller *Controller) ForgotPassword(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[dto.ForgotPasswordRequest, dto.AcceptedResult]{
		Decode: httprequest.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.ForgotPasswordRequest) (rest.Result[dto.AcceptedResult], error) {
			// A public recovery request must not reveal account eligibility.
			_ = controller.service.Forgot(ctx, input.RequestBody.Email)
			return rest.Result[dto.AcceptedResult]{
				StatusCode: http.StatusAccepted,
				Data: dto.AcceptedResult{Accepted: true, Delivery: "EMAIL"},
			}, nil
		},
	})
}

func (controller *Controller) exchange(writer http.ResponseWriter, request *http.Request, prefix string, purpose string, cookie string, path string) {
	rest.RunStruct(writer, request, rest.StructConfig[dto.EmptyRequest, dto.GrantResult]{
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, _ dto.EmptyRequest) (rest.Result[dto.GrantResult], error) {
			grant, err := controller.service.Exchange(ctx, controller.cookies.Proof(request, prefix), purpose)
			if err != nil {
				return rest.Result[dto.GrantResult]{}, helpers.Map(err)
			}
			csrf, err := secrets.NewSecret(32)
			if err != nil {
				return rest.Result[dto.GrantResult]{}, helpers.Map(err)
			}
			controller.cookies.SetGrant(writer, cookie, path, grant.Secret)
			http.SetCookie(writer, &http.Cookie{
				Name: cookies.CSRFName,
				Value: csrf,
				Path: "/",
				MaxAge: 600,
				Secure: true,
				SameSite: http.SameSiteStrictMode,
			})
			return rest.OK(grant.Data), nil
		},
	})
}


func (controller *Controller) inspectProof(writer http.ResponseWriter, request *http.Request, prefix string, purpose string) {
  writer.Header().Set("Cache-Control", "no-store")
  rest.RunStruct(writer, request, rest.StructConfig[dto.EmptyRequest, dto.GrantResult]{
    Fallback: definitions.Internal,
    Service: func(ctx context.Context, _ dto.EmptyRequest) (rest.Result[dto.GrantResult], error) {
      result, err := controller.service.InspectProof(ctx, controller.cookies.Proof(request, prefix), purpose)
      if err != nil { return rest.Result[dto.GrantResult]{}, helpers.Map(err) }
      return rest.OK(result), nil
    },
  })
}

func (controller *Controller) RecoveryValidate(writer http.ResponseWriter, request *http.Request) {
  controller.inspectProof(writer, request, "Recovery", "recovery")
}

func (controller *Controller) InvitationValidate(writer http.ResponseWriter, request *http.Request) {
  controller.inspectProof(writer, request, "Invitation", "invitation")
}

func (controller *Controller) RecoveryExchange(writer http.ResponseWriter, request *http.Request) {
	controller.exchange(writer, request, "Recovery", "recovery", cookies.ResetCookie, "/identity/public/auth/reset-password")
}

func (controller *Controller) InvitationExchange(writer http.ResponseWriter, request *http.Request) {
	controller.exchange(writer, request, "Invitation", "invitation", cookies.SetupCookie, "/identity/public/accounts/setup")
}

func (controller *Controller) ResetPassword(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[dto.ResetPasswordRequest, dto.StatusResult]{
		Decode: httprequest.DecodeStrictJSON,
		InvalidRequestBody: definitions.InvalidRequest,
		Fallback: definitions.Internal,
		Service: func(ctx context.Context, input dto.ResetPasswordRequest) (rest.Result[dto.StatusResult], error) {
			if err := controller.service.ResetPassword(ctx, controller.cookies.Read(request, cookies.ResetCookie), input.RequestBody.NewPassword); err != nil {
				return rest.Result[dto.StatusResult]{}, helpers.Map(err)
			}
			controller.cookies.ClearGrant(writer, cookies.ResetCookie, "/identity/public/auth/reset-password")
			return rest.OK(dto.StatusResult{Status: "PASSWORD_RESET", SessionsRevoked: true}), nil
		},
	})
}

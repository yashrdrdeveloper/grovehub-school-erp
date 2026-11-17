package impl

import (
	"context"
	"net/http"

	httprequest "grovehub/libs/core/http/request"
	rest "grovehub/libs/core/pipeline/rest"
	errorcode "grovehub/services/identity/internal/app/constants/error"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	authdto "grovehub/services/identity/internal/app/dto/authentication"
	protected "grovehub/services/identity/internal/app/controller/authentication/protected"
	"grovehub/services/identity/internal/app/controller/helpers"
	auth "grovehub/services/identity/internal/app/services/authentication"
	"grovehub/services/identity/internal/app/transport/cookies"
)

type Controller struct {
	cookies cookies.Manager
	service auth.Service
}

var _ protected.AuthenticationController = (*Controller)(nil)

func New(service auth.Service, manager cookies.Manager) protected.AuthenticationController {
	return &Controller{service: service, cookies: manager}
}

func (controller *Controller) Logout(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[authdto.EmptyRequest, authdto.StatusResult]{
		Fallback: errorcode.Internal,
		Service: func(ctx context.Context, _ authdto.EmptyRequest) (rest.Result[authdto.StatusResult], error) {
			if err := controller.service.Logout(ctx); err != nil {
				return rest.Result[authdto.StatusResult]{}, helpers.Map(err)
			}
			controller.cookies.ClearSession(writer)
			return rest.OK(authdto.StatusResult{Status: "LOGGED_OUT"}), nil
		},
	})
}

func (controller *Controller) LogoutAll(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[authdto.EmptyRequest, authdto.RevocationResult]{
		Fallback: errorcode.Internal,
		Service: func(ctx context.Context, _ authdto.EmptyRequest) (rest.Result[authdto.RevocationResult], error) {
			count, err := controller.service.LogoutAll(ctx)
			if err != nil {
				return rest.Result[authdto.RevocationResult]{}, helpers.Map(err)
			}
			controller.cookies.ClearSession(writer)
			return rest.OK(authdto.RevocationResult{
				Status: "ALL_SESSIONS_REVOKED",
				RevokedCount: count,
			}), nil
		},
	})
}

func (controller *Controller) ChangePassword(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[authdto.ChangePasswordRequest, authdto.StatusResult]{
		Decode: httprequest.DecodeStrictJSON,
		InvalidRequestBody: errorcode.InvalidRequest,
		Fallback: errorcode.Internal,
		Service: func(ctx context.Context, input authdto.ChangePasswordRequest) (rest.Result[authdto.StatusResult], error) {
			if err := controller.service.ChangePassword(ctx, input.RequestBody.NewPassword); err != nil {
				return rest.Result[authdto.StatusResult]{}, helpers.Map(err)
			}
			controller.cookies.ClearSession(writer)
			return rest.OK(authdto.StatusResult{
				Status: "PASSWORD_CHANGED",
				SessionsRevoked: true,
			}), nil
		},
	})
}

func (controller *Controller) Reauthenticate(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[authdto.ReauthenticateRequest, authdto.StatusResult]{
		Decode: httprequest.DecodeStrictJSON,
		InvalidRequestBody: errorcode.InvalidRequest,
		Fallback: errorcode.Internal,
		Service: func(ctx context.Context, input authdto.ReauthenticateRequest) (rest.Result[authdto.StatusResult], error) {
			if err := controller.service.Reauthenticate(ctx, input.RequestBody.Password); err != nil {
				return rest.Result[authdto.StatusResult]{}, helpers.Map(err)
			}
			return rest.OK(authdto.StatusResult{Status: "REAUTHENTICATED"}), nil
		},
	})
}

func (controller *Controller) Sessions(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[authdto.EmptyRequest, authdto.SessionsResult]{
		Fallback: errorcode.Internal,
		Service: func(ctx context.Context, _ authdto.EmptyRequest) (rest.Result[authdto.SessionsResult], error) {
			result, err := controller.service.ListSessions(ctx)
			return rest.OK(result), helpers.Map(err)
		},
	})
}

func (controller *Controller) RevokeSession(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[authdto.EmptyRequest, authdto.SessionRevokedResult]{
		Fallback: errorcode.Internal,
		Service: func(ctx context.Context, _ authdto.EmptyRequest) (rest.Result[authdto.SessionRevokedResult], error) {
			sessionID := request.PathValue("sessionId")
			if err := controller.service.RevokeSession(ctx, sessionID); err != nil {
				return rest.Result[authdto.SessionRevokedResult]{}, helpers.Map(err)
			}
			return rest.OK(authdto.SessionRevokedResult{
				SessionID: sessionID,
				Status: "REVOKED",
			}), nil
		},
	})
}

func (controller *Controller) Me(writer http.ResponseWriter, request *http.Request) {
	rest.RunStruct(writer, request, rest.StructConfig[authdto.EmptyRequest, accountdto.Account]{
		Fallback: errorcode.Internal,
		Service: func(ctx context.Context, _ authdto.EmptyRequest) (rest.Result[accountdto.Account], error) {
			account, err := controller.service.Me(ctx)
			return rest.OK(account), helpers.Map(err)
		},
	})
}

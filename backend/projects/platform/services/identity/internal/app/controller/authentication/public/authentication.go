package public

import "net/http"

// @Routes /auth
type AuthenticationController interface {
	// @POST /login
	// @RouteLimit 12
	Login(writer http.ResponseWriter, request *http.Request)

	// @POST /refresh
	// @RouteLimit 30
	Refresh(writer http.ResponseWriter, request *http.Request)

	// @POST /forgot-password
	// @RouteLimit 5
	ForgotPassword(writer http.ResponseWriter, request *http.Request)

	// @GET /recovery/validate
	RecoveryValidate(writer http.ResponseWriter, request *http.Request)

	// @GET /invitations/validate
	InvitationValidate(writer http.ResponseWriter, request *http.Request)

	// @POST /recovery/exchange
	// @RouteLimit 8
	RecoveryExchange(writer http.ResponseWriter, request *http.Request)

	// @GET /reset-password/status
	ResetGrantStatus(writer http.ResponseWriter, request *http.Request)

	// @POST /reset-password
	ResetPassword(writer http.ResponseWriter, request *http.Request)

	// @POST /invitations/exchange
	// @RouteLimit 8
	InvitationExchange(writer http.ResponseWriter, request *http.Request)
}

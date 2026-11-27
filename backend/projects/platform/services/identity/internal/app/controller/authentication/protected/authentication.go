package protected

import "net/http"

// @Routes /auth
type AuthenticationController interface {
	// @POST /logout
	Logout(writer http.ResponseWriter, request *http.Request)

	// @POST /logout-all
	LogoutAll(writer http.ResponseWriter, request *http.Request)

	// @POST /change-password
	ChangePassword(writer http.ResponseWriter, request *http.Request)

	// @POST /reauthenticate
	Reauthenticate(writer http.ResponseWriter, request *http.Request)

	// @GET /sessions
	Sessions(writer http.ResponseWriter, request *http.Request)

	// @DELETE /sessions/{sessionId}
	RevokeSession(writer http.ResponseWriter, request *http.Request)

	// @GET /me
	Me(writer http.ResponseWriter, request *http.Request)
}

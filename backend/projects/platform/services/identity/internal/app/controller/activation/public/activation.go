package public

import "net/http"

// @Routes /accounts
// Invitation setup is proof/grant protected, not open account registration.
type ActivationController interface {
	// @GET /setup/status
	SetupStatus(writer http.ResponseWriter, request *http.Request)

	// @POST /setup
	Setup(writer http.ResponseWriter, request *http.Request)
}

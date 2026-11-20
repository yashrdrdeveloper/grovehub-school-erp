package public

import "net/http"

// @Routes /accounts
type AccountController interface {
	// @GET /setup/status
	SetupStatus(writer http.ResponseWriter, request *http.Request)

	// @POST /setup
	Setup(writer http.ResponseWriter, request *http.Request)
}

package public

import "net/http"

// @Routes /accounts
type AccountController interface {
	// @POST /setup
	Setup(writer http.ResponseWriter, request *http.Request)
}

package public

import "net/http"

// OAuth callbacks validate their own state; an access session is not required.
// @Routes /connections
type ConnectionController interface {
	// @GET /oauth/{provider}/callback
	FinishOAuth(http.ResponseWriter, *http.Request)
}

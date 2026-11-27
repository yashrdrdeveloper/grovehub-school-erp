package public

import "net/http"

// @Routes /providers
type ProviderController interface {
	// @GET /
	Providers(http.ResponseWriter, *http.Request)
	// @GET /{provider}
	Provider(http.ResponseWriter, *http.Request)
}

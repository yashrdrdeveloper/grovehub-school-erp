package protected

import "net/http"

// @Routes /capabilities
type CapabilityController interface {
	// @GET /me
	Discover(writer http.ResponseWriter, request *http.Request)

	// @GET /permissions
	Catalog(writer http.ResponseWriter, request *http.Request)
}

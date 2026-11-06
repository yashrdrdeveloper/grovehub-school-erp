package routes

import (
	"net/http"

	"grovehub/services/integration/controller"
)

// RegisterHTTP composes routes by URL resource, not internal capability.
func RegisterHTTP(mux *http.ServeMux, c controller.Controller) {
	registerProviderRoutes(mux, c)
	registerConnectionRoutes(mux, c)
}

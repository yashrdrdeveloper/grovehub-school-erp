package helpers

import (
	"net/http"

	httptracing "grovehub/libs/core/http/tracing"
)

// Handle registers a traced HTTP endpoint for the Integration service.
func Handle(mux *http.ServeMux, method, path string, handler http.HandlerFunc) {
	pattern := method + " " + path
	mux.Handle(pattern, httptracing.ServiceHandler("integration", pattern, handler))
}

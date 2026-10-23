package routes

import (
	"net/http"

	httptracing "grovehub/libs/core/http/tracing"
)

func RegisterHTTP(mux *http.ServeMux, healthController http.Handler) {
	mux.Handle(
		"GET /identity/health",
		httptracing.ServiceHandler("identity", "GET /identity/health", healthController),
	)
}

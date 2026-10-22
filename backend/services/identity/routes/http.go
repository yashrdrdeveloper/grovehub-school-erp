package routes

import "net/http"

func RegisterHTTP(mux *http.ServeMux, healthController http.Handler) {
	mux.Handle("GET /identity/health", healthController)
}

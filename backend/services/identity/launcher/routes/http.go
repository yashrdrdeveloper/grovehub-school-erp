package routes

import "net/http"

func RegisterHTTP(mux *http.ServeMux, healthHandler http.Handler) {
	mux.Handle("GET /identity/health", healthHandler)
}

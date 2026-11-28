package routing

import (
	"net/http"

	httpresponse "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http/response"
)

type healthResponse struct {
	Status string `json:"status"`
}

func RegisterHealth(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", func(writer http.ResponseWriter, _ *http.Request) {
		_ = httpresponse.WriteSuccess(writer, http.StatusOK, healthResponse{Status: "UP"})
	})
}

// Readiness intentionally returns only aggregate status, never dependency or
// tenant details. Liveness remains independent of dependencies.
func RegisterReadiness(mux *http.ServeMux,ready func() bool) {
	mux.HandleFunc("GET /ready",func(writer http.ResponseWriter,_ *http.Request) {
		writer.Header().Set("Cache-Control","no-store")
		if ready==nil || !ready() {
			_ = httpresponse.WriteSuccess(writer,http.StatusServiceUnavailable,healthResponse{Status:"NOT_READY"})
			return
		}
		_ = httpresponse.WriteSuccess(writer,http.StatusOK,healthResponse{Status:"READY"})
	})
}

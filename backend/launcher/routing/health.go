package routing

import (
	"net/http"

	httpresponse "grovehub/libs/core/http/response"
)

type healthResponse struct {
	Status string `json:"status"`
}

func RegisterHealth(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", func(writer http.ResponseWriter, _ *http.Request) {
		_ = httpresponse.WriteSuccess(writer, http.StatusOK, healthResponse{Status: "UP"})
	})
}

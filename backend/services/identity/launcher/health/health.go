package health

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(writer).Encode(response{
		Service: "identity",
		Status:  "ok",
	})
}

package controller

import (
	"encoding/json"
	"net/http"
)

type healthResponse struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (controller *HealthController) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(writer).Encode(healthResponse{
		Service: "identity",
		Status:  "ok",
	})
}

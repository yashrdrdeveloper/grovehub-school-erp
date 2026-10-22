package controller

import (
	"net/http"

	httpresponse "grovehub/libs/core/http/response"
)

type healthData struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (controller *HealthController) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	_ = httpresponse.WriteSuccess(writer, http.StatusOK, healthData{
		Service: "identity",
		Status:  "ok",
	})
}

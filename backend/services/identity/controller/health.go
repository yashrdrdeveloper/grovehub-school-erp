package controller

import (
	"net/http"

	httpresponse "grovehub/libs/core/http/response"
	identityresponse "grovehub/services/identity/dto/v1/response"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (controller *HealthController) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	_ = httpresponse.WriteSuccess(writer, http.StatusOK, identityresponse.HealthData{
		Service: "identity",
		Status:  "ok",
	})
}

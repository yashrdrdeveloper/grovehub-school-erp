package routes

import (
	"net/http"

	httptracing "grovehub/libs/core/http/tracing"
	"grovehub/services/identity/controller"
)

func RegisterHTTP(mux *http.ServeMux, registrationController controller.RegistrationController) {
	mux.Handle(
		"POST /auth/register",
		httptracing.ServiceHandler(
			"identity",
			"POST /auth/register",
			http.HandlerFunc(registrationController.Register),
		),
	)
}

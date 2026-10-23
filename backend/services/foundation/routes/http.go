package routes

import (
	"net/http"

	httptracing "grovehub/libs/core/http/tracing"
	"grovehub/services/foundation/controller"
)

func RegisterHTTP(mux *http.ServeMux, personController *controller.PersonController) {
	mux.Handle(
		"POST /foundation/persons",
		httptracing.ServiceHandler(
			"foundation",
			"POST /foundation/persons",
			http.HandlerFunc(personController.Create),
		),
	)

	mux.Handle(
		"GET /foundation/persons/{gid}",
		httptracing.ServiceHandler(
			"foundation",
			"GET /foundation/persons/{gid}",
			http.HandlerFunc(personController.GetByGID),
		),
	)

	mux.Handle(
		"PATCH /foundation/persons/{gid}/deactivate",
		httptracing.ServiceHandler(
			"foundation",
			"PATCH /foundation/persons/{gid}/deactivate",
			http.HandlerFunc(personController.Deactivate),
		),
	)
}

package routes

import (
	"net/http"

	"grovehub/launcher/internal/src/controller"
)

func RegisterHTTP(mux *http.ServeMux, health controller.HealthController) {
	registerHealthRoutes(mux, health)
}

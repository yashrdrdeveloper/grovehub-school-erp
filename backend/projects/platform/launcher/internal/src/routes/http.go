package routes

import (
	"grovehub/launcher/internal/src/controller"
	"net/http"
)

func RegisterHTTP(mux *http.ServeMux, health controller.HealthController) {
	registerHealthRoutes(mux, health)
}

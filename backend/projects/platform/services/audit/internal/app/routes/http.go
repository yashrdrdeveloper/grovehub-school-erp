package routes

import (
	protected "grovehub/services/audit/internal/app/controller/activity/protected"
	"net/http"
)

// @Routes /audit
func RegisterHTTP(mux *http.ServeMux, controller protected.ActivityController) {
	registerProtectedActivityRoutes(mux, controller)
}

package routes

import (
	"net/http"
	protected "grovehub/services/audit/internal/app/controller/activity/protected"
)

// @Routes /audit
func RegisterHTTP(mux *http.ServeMux, controller protected.ActivityController) {
	registerProtectedActivityRoutes(mux,controller)
}

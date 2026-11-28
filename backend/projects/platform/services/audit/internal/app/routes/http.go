package routes

import (
	"net/http"

	protected "grovehub/services/audit/internal/app/controller/activity/protected"
	restricted "grovehub/services/audit/internal/app/controller/operator/restricted"
)

// @Routes /audit
func RegisterHTTP(mux *http.ServeMux, activityController protected.ActivityController, operatorController restricted.OperatorController) {
	registerProtectedActivityRoutes(mux, activityController)
	registerRestrictedOperatorRoutes(mux, operatorController)
}

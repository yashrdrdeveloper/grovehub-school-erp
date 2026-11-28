package routes_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	activitycontroller "grovehub/services/audit/internal/app/controller/activity/protected/impl"
	"grovehub/services/audit/internal/app/controller/operator/constants"
	operatorcontroller "grovehub/services/audit/internal/app/controller/operator/restricted/impl"
	"grovehub/services/audit/internal/app/routes"
)

func TestGeneratedRestrictedRoutes(t *testing.T) {
	mux := http.NewServeMux()
	routes.RegisterHTTP(mux, activitycontroller.New(nil), operatorcontroller.New(nil))

	for _, pattern := range []string{
		constants.Pattern,
		constants.SearchPattern,
		constants.SchoolCheckPattern,
		constants.ReceiptPattern,
		constants.StatusReceiptPattern,
		constants.ReceiptSearchPattern,
	} {
		method, path, _ := strings.Cut(pattern, " ")
		path = strings.ReplaceAll(path, "{schoolId}", "e93696cc-15f0-4e06-9577-0d5cfc2a83c8")
		t.Run(pattern, func(t *testing.T) {
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, httptest.NewRequest(method, path, nil))
			if w.Code != http.StatusForbidden || w.Body.Len() != 0 {
				t.Fatalf("route missing or unauthenticated access permitted: status=%d", w.Code)
			}
		})
	}
}

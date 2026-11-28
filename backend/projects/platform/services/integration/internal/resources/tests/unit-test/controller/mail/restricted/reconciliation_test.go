package restricted_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"grovehub/services/integration/internal/app/controller/mail/restricted"
)
func TestReconciliationRejectsAnonymousRequests(t *testing.T){
	w:=httptest.NewRecorder()
	r:=httptest.NewRequest(http.MethodPost,"/integrations/restricted/mail/reconciliation/search",strings.NewReader(`{"tenantId":"school","limit":10}`))
	restricted.SearchHandler(nil).ServeHTTP(w,r)
	if w.Code!=http.StatusForbidden||w.Header().Get("Cache-Control")!="no-store"{t.Fatal(w.Code,w.Header())}
}

package correlation_test

import (
	corehttp "grovehub/libs/core/http"

	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestIngress(t *testing.T) {
	h := corehttp.Ingress(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := corehttp.ID(r.Context())
		if _, err := uuid.Parse(id); err != nil || w.Header().Get(corehttp.RequestID) != id {
			t.Fatal("invalid request correlation")
		}
		if corehttp.ID(corehttp.Ensure(r.Context())) != id {
			t.Fatal("correlation changed")
		}
	}))
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set(corehttp.RequestID, "forged")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Header().Get(corehttp.RequestID) == "forged" || corehttp.ID(corehttp.WithNew(context.Background())) == w.Header().Get(corehttp.RequestID) {
		t.Fatal("correlation reused")
	}
}

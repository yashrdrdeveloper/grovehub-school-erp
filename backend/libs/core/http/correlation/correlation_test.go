package correlation

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestIngress(t *testing.T) {
	h := Ingress(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := ID(r.Context())
		if _, err := uuid.Parse(id); err != nil || w.Header().Get(Header) != id {
			t.Fatal("invalid request correlation")
		}
		if ID(Ensure(r.Context())) != id {
			t.Fatal("correlation changed")
		}
	}))
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set(Header, "forged")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Header().Get(Header) == "forged" || ID(WithNew(context.Background())) == w.Header().Get(Header) {
		t.Fatal("correlation reused")
	}
}

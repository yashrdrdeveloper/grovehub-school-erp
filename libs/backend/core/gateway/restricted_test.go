package gateway

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRestrictedOnly(t *testing.T) {
	mux := http.NewServeMux()
	for _, pattern := range []string{"GET /identity/restricted/accounts", "GET /identity/public/accounts", "GET /identity/protected/accounts", "GET /health", "GET /identity/public/restricted/accounts"} {
		mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})
	}
	handler := RestrictedOnly(mux, mux)
	cases := []struct {
		method string
		path string
		want int
	}{
		{http.MethodGet, "/identity/restricted/accounts", http.StatusNoContent},
		{http.MethodGet, "/identity/public/accounts", http.StatusForbidden},
		{http.MethodGet, "/identity/protected/accounts", http.StatusForbidden},
		{http.MethodGet, "/health", http.StatusForbidden},
		{http.MethodGet, "/identity/public/restricted/accounts", http.StatusForbidden},
		{http.MethodGet, "/identity/restricted/missing", http.StatusForbidden},
		{http.MethodPost, "/identity/restricted/accounts", http.StatusForbidden},
	}
	for _, tc := range cases {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, nil))
			if w.Code != tc.want {
				t.Fatalf("status = %d, want %d", w.Code, tc.want)
			}
		})
	}
}

func TestRestrictedOnlyMissingDependencies(t *testing.T) {
	w := httptest.NewRecorder()
	RestrictedOnly(nil, nil).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/identity/restricted/accounts", nil))
	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want forbidden", w.Code)
	}
}

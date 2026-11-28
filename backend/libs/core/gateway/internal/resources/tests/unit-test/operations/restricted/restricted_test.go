package restricted_test

import (
	"grovehub/libs/core/gateway/internal/src/operations/restricted"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RestrictedOnly", func() {
	It("dispatches only registered restricted routes and fails closed otherwise", func() {
		mux := http.NewServeMux()
		for _, path := range []string{"/identity/restricted/accounts", "/identity/public/accounts", "/identity/protected/accounts", "/health"} {
			mux.HandleFunc("GET "+path, func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })
		}
		handler := restricted.RestrictedOnly(mux, mux)
		for _, path := range []string{"/identity/public/accounts", "/identity/protected/accounts", "/health", "/identity/restricted/missing"} {
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
			Expect(w.Code).To(Equal(http.StatusForbidden))
			Expect(w.Header().Get("Cache-Control")).To(Equal("no-store"))
		}
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/identity/restricted/accounts", nil))
		Expect(w.Code).To(Equal(http.StatusNoContent))
		for _, handler := range []http.Handler{restricted.RestrictedOnly(nil, mux), restricted.RestrictedOnly(mux, nil)} {
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/identity/restricted/accounts", nil))
			Expect(w.Code).To(Equal(http.StatusForbidden))
		}
	})
})

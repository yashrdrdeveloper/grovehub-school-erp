package readiness_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"grovehub/libs/core/http/internal/src/operations/readiness"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Handler", func() {
	It("fails closed when draining, missing checks or a dependency fails", func() {
		var drain atomic.Bool
		ok := func(context.Context) error { return nil }
		failed := func(context.Context) error { return errors.New("private detail") }
		for _, entry := range []struct { handler http.HandlerFunc; status int }{
			{readiness.Handler(nil, ok),503}, {readiness.Handler(&drain),503},
			{readiness.Handler(&drain,nil),503}, {readiness.Handler(&drain,failed),503},
			{readiness.Handler(&drain,ok),204},
		} {
			w := httptest.NewRecorder()
			entry.handler.ServeHTTP(w, httptest.NewRequest("GET","/ready",nil))
			Expect(w.Code).To(Equal(entry.status))
			Expect(w.Body.String()).To(BeEmpty())
			Expect(w.Header().Get("Cache-Control")).To(Equal("no-store"))
		}
		drain.Store(true)
		w := httptest.NewRecorder()
		readiness.Handler(&drain,ok).ServeHTTP(w, httptest.NewRequest("GET","/ready",nil))
		Expect(w.Code).To(Equal(503))
	})
})

var _ = Describe("HTTP", func() {
	It("allows only fixed safe origins and readiness paths", func() {
		for _, origin := range []string{"http://example.com","https://user@example.com","https://example.com/path","https://example.com?","https://example.com#x"} {
			_, err := readiness.HTTP(origin,"/ready")
			Expect(err).To(HaveOccurred())
		}
		_, err := readiness.HTTP("https://example.com","/other")
		Expect(err).To(HaveOccurred())
		_, err = readiness.HTTP("https://example.com","/ready")
		Expect(err).NotTo(HaveOccurred())
	})
})

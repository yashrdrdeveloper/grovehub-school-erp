package readiness_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/http/internal/src/operations/readiness"
)

var _ = Describe("Readiness probe edge cases", func() {
	It("rejects a dependency that begins draining during checks", func() {
		var drain atomic.Bool
		w := httptest.NewRecorder()
		readiness.Handler(&drain, func(context.Context) error { drain.Store(true); return nil })(w, httptest.NewRequest("GET", "/ready", nil))
		Expect(w.Code).To(Equal(503))
	})
	It("checks loopback status, rejects redirects and handles cancellation", func() {
		var status atomic.Int32
		status.Store(200)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/ready" {
				w.WriteHeader(404)
				return
			}
			if status.Load() == 302 {
				w.Header().Set("Location", "/ready")
			}
			w.WriteHeader(int(status.Load()))
		}))
		defer server.Close()
		check, err := readiness.HTTP(server.URL, "/ready")
		Expect(err).NotTo(HaveOccurred())
		Expect(check(context.Background())).To(Succeed())
		status.Store(204)
		Expect(check(context.Background())).To(Succeed())
		for _, code := range []int32{302, 503} {
			status.Store(code)
			Expect(check(context.Background())).To(MatchError("dependency not ready"))
		}
		Expect(check(nil)).To(MatchError("dependency probe failed"))
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		Expect(check(ctx)).To(MatchError("dependency unavailable"))
	})
})

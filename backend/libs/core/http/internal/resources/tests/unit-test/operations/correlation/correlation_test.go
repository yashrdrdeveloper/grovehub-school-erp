package correlation_test

import (
	"context"
	"github.com/google/uuid"
	"grovehub/libs/core/http/internal/src/operations/correlation"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ingress", func() {
	It("generates server IDs, ignores inbound IDs and preserves an existing context ID only when ensuring", func() {
		Expect(correlation.ID(context.Background())).To(BeEmpty())
		ctx := correlation.Ensure(context.Background())
		_, err := uuid.Parse(correlation.ID(ctx))
		Expect(err).NotTo(HaveOccurred())
		Expect(correlation.ID(correlation.Ensure(ctx))).To(Equal(correlation.ID(ctx)))
		h := correlation.Ingress(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Expect(correlation.ID(r.Context())).NotTo(Equal("forged"))
			Expect(w.Header().Get(correlation.Header)).To(Equal(correlation.ID(r.Context())))
		}))
		r := httptest.NewRequest("GET", "/", nil)
		r.Header.Set(correlation.Header, "forged")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		Expect(w.Header().Get(correlation.Header)).NotTo(BeEmpty())
	})
})

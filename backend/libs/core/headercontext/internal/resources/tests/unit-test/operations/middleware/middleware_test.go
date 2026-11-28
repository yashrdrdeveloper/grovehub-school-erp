package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"grovehub/libs/core/headercontext/internal/src/contracts"
	"grovehub/libs/core/headercontext/internal/src/operations/metadata"
	"grovehub/libs/core/headercontext/internal/src/operations/middleware"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("New", func() {
	It("attaches validated values without changing the incoming request", func() {
		parent := metadata.WithValues(context.Background(), contracts.Values{IdempotencyKey: "parent"})
		request := httptest.NewRequest(http.MethodPost, "/", nil).WithContext(parent)
		request.Header.Set("Idempotency-Key", "child")
		request.Header.Set("X-Correlation-ID", "corr")
		calls := 0
		handler := middleware.New(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			calls++
			Expect(metadata.FromContext(r.Context())).To(Equal(contracts.Values{IdempotencyKey: "child", CorrelationID: "corr"}))
			w.WriteHeader(http.StatusNoContent)
		}))
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		Expect(calls).To(Equal(1))
		Expect(recorder.Code).To(Equal(http.StatusNoContent))
		Expect(metadata.FromContext(request.Context()).IdempotencyKey).To(Equal("parent"))
	})

	It("returns a generic bad request and never invokes downstream on invalid metadata", func() {
		handler := middleware.New(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			Fail("invalid metadata reached downstream")
		}))
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.Header.Set("X-Request-ID", "private invalid value")
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusBadRequest))
		Expect(recorder.Body.String()).To(Equal("invalid request metadata header\n"))
	})
})

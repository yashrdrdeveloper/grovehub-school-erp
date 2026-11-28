package headercontext_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coreevents "grovehub/libs/core/events"
	headercontext "grovehub/libs/core/headercontext"
)

var _ = Describe("Header context", func() {
	It("propagates validated metadata and event correlation", func() {
		handler := headercontext.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			Expect(headercontext.IdempotencyKey(r.Context())).To(Equal("create:123"))
			Expect(headercontext.CorrelationID(r.Context())).To(Equal("correlation-123"))
			Expect(headercontext.RequestID(r.Context())).To(Equal("request-123"))
			event := coreevents.NewEnvelope(r.Context(), coreevents.Metadata{EventID:"event", EventType:"test", Version:1}, "payload")
			Expect(event.CorrelationID).To(Equal("correlation-123"))
			w.WriteHeader(http.StatusNoContent)
		}))
		request := httptest.NewRequest(http.MethodPost, "/", nil)
		request.Header.Set("Idempotency-Key", " create:123 ")
		request.Header.Set("X-Correlation-ID", "correlation-123")
		request.Header.Set("X-Request-ID", "request-123")
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		Expect(recorder.Code).To(Equal(http.StatusNoContent))
	})

	It("rejects duplicate, blank, and invalid metadata before invoking business logic", func() {
		handler := headercontext.Middleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { Fail("invalid request must not reach handler") }))
		for _, values := range [][]string{{"valid", "other"}, {"invalid key"}, {" "}} {
			request := httptest.NewRequest(http.MethodPost, "/", nil)
			for _, value := range values { request.Header.Add("Idempotency-Key", value) }
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			Expect(recorder.Body.String()).NotTo(ContainSubstring("invalid key"))
		}
	})

	It("does not trust incoming trace identifiers or leak a parent idempotency key", func() {
		handler := headercontext.Middleware(http.HandlerFunc(func(_ http.ResponseWriter, request *http.Request) {
			Expect(headercontext.FromContext(request.Context()).CorrelationID).To(BeEmpty())
			Expect(headercontext.FromContext(request.Context()).IdempotencyKey).To(Equal("register-key"))
		}))
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.Header.Set("Idempotency-Key", "register-key")
		request.Header.Set("X-Trace-ID", "spoofed")
		handler.ServeHTTP(httptest.NewRecorder(), request)
		ctx := headercontext.WithValues(context.Background(), headercontext.Values{IdempotencyKey:"register-key", CorrelationID:"corr"})
		internal := headercontext.WithoutIdempotencyKey(ctx)
		Expect(headercontext.IdempotencyKey(internal)).To(BeEmpty())
		Expect(headercontext.CorrelationID(internal)).To(Equal("corr"))
		Expect(headercontext.IdempotencyKey(ctx)).To(Equal("register-key"))
	})
})

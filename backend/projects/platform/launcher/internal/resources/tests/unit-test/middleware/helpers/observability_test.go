package helpers_test

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"grovehub/launcher/internal/src/middleware/helpers"
)

// @Testing suite=LauncherResponseRecorder
func ObservabilitySpecs() {
	Describe("Observability helpers", func() {
		It("records success and failure spans with the appropriate status and event", func() {
			spanRecorder := tracetest.NewSpanRecorder()
			provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(spanRecorder))
			DeferCleanup(func() {
				Expect(provider.Shutdown(context.Background())).To(Succeed())
			})
			request := httptest.NewRequest(http.MethodGet, "/health", nil)
			for _, status := range []int{http.StatusOK, http.StatusBadRequest} {
				ctx, span := provider.Tracer("launcher-test").Start(context.Background(), "request")
				recorder := &helpers.Recorder{ResponseWriter: httptest.NewRecorder(), Status: status}
				helpers.RecordTrace(ctx, request, "/health", recorder, time.Millisecond)
				span.End()
			}

			spans := spanRecorder.Ended()
			Expect(spans).To(HaveLen(2))
			Expect(spans[0].Status().Code).To(Equal(codes.Unset))
			Expect(spans[0].Name()).To(HaveSuffix("GET /health"))
			Expect(spans[0].Events()).To(BeEmpty())
			Expect(spans[1].Status().Code).To(Equal(codes.Error))
			Expect(spans[1].Events()).To(HaveLen(1))
			Expect(spans[1].Events()[0].Name).To(Equal("http.request.failed"))
		})
		It("records the HTTP status and duration without requiring a configured metrics exporter", func() {
			var buffer bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&buffer, nil))
			request := httptest.NewRequest(http.MethodGet, "/health", nil)
			recorder := &helpers.Recorder{ResponseWriter: httptest.NewRecorder(), Status: http.StatusOK, Bytes: 7}
			helpers.RecordLog(logger, request.Context(), request, "/health", recorder, 5*time.Millisecond)
			Expect(buffer.String()).To(ContainSubstring("request completed"))
			Expect(buffer.String()).To(ContainSubstring("\"status\":200"))
			Expect(buffer.String()).To(ContainSubstring("\"durationMs\":5"))
			helpers.RecordMetrics(request.Context(), request, "/health", recorder, 5*time.Millisecond)
		})
	})
}

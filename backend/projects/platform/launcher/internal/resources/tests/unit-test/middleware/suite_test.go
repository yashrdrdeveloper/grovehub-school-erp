package middleware_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	oteltrace "go.opentelemetry.io/otel/trace"

	coretelemetry "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/telemetry"
	launcherMiddleware "grovehub/launcher/internal/src/middleware/observability/impl"
)

const failedPath = "/abc"

type observedRequest struct {
	recorder *httptest.ResponseRecorder
	log      map[string]any
	traceID  string
	spanID   string
}

var (
	testTracerProvider *sdktrace.TracerProvider
	spanRecorder       *tracetest.SpanRecorder
)

// @BeforeEach
func setupObservability() {
	spanRecorder = tracetest.NewSpanRecorder()
	testTracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(spanRecorder),
	)
	otel.SetTracerProvider(testTracerProvider)
}

var _ = AfterEach(func() {
	Expect(testTracerProvider.Shutdown(context.Background())).To(Succeed())
	otel.SetTracerProvider(oteltrace.NewNoopTracerProvider())
})

func serveObservedRequest(status int) observedRequest {
	ctx, span := coretelemetry.StartSpan(context.Background(), "edge", "GET "+failedPath)
	traceID := coretelemetry.TraceID(ctx)
	spanID := coretelemetry.SpanID(ctx)

	buffer := &bytes.Buffer{}
	logger := slog.New(slog.NewJSONHandler(buffer, nil))
	handler := launcherMiddleware.New(logger).Wrap(
		http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
			writer.WriteHeader(status)
			_, _ = writer.Write([]byte("response"))
		}),
	)

	request := httptest.NewRequest(http.MethodGet, failedPath, nil).WithContext(ctx)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	span.End()

	var logPayload map[string]any
	Expect(json.Unmarshal(buffer.Bytes(), &logPayload)).To(Succeed())

	return observedRequest{
		recorder: recorder,
		log:      logPayload,
		traceID:  traceID,
		spanID:   spanID,
	}
}

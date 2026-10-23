package tracing_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"

	httpheaders "grovehub/libs/core/http/headers"
	httptracing "grovehub/libs/core/http/tracing"
	coretelemetry "grovehub/libs/core/telemetry"
)

const (
	identityHealthPath       = "/identity/health"
	identityService          = "identity"
	identityOperation        = "GET /identity/health"
	internalServiceURL       = "http://service.test/internal"
	backgroundServiceURL     = "http://service.test/background"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (roundTrip roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return roundTrip(request)
}

var testTracerProvider *sdktrace.TracerProvider

func TestCoreHTTPTracingUnit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Core HTTP Tracing Unit Tests")
}

var _ = BeforeEach(func() {
	testTracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	otel.SetTracerProvider(testTracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
})

var _ = AfterEach(func() {
	Expect(testTracerProvider.Shutdown(context.Background())).To(Succeed())
	otel.SetTracerProvider(oteltrace.NewNoopTracerProvider())
})

func newTracedInboundRequest(method string, path string) (*http.Request, context.Context, oteltrace.Span) {
	parentContext, parentSpan := coretelemetry.StartSpan(context.Background(), "edge", method+" "+path)
	request := httptest.NewRequest(method, path, nil).WithContext(parentContext)
	return request, parentContext, parentSpan
}

func newTracedOutboundRequest(method string, targetURL string) (*http.Request, context.Context, oteltrace.Span) {
	parentContext, parentSpan := coretelemetry.StartSpan(context.Background(), "edge", method+" "+targetURL)
	request, err := http.NewRequestWithContext(parentContext, method, targetURL, nil)
	Expect(err).NotTo(HaveOccurred())
	return request, parentContext, parentSpan
}

func newOutboundRequest(method string, targetURL string) *http.Request {
	request, err := http.NewRequest(method, targetURL, nil)
	Expect(err).NotTo(HaveOccurred())
	return request
}

func noContentResponse(request *http.Request) *http.Response {
	return &http.Response{
		StatusCode: http.StatusNoContent,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader("")),
		Request:    request,
	}
}

func captureOutboundTraceHeaders(request *http.Request) (string, string) {
	var capturedTraceID string
	var capturedTraceParent string
	client := &http.Client{
		Transport: httptracing.NewTransport(roundTripFunc(func(outbound *http.Request) (*http.Response, error) {
			capturedTraceID = outbound.Header.Get(httpheaders.TraceID)
			capturedTraceParent = outbound.Header.Get("traceparent")
			return noContentResponse(outbound), nil
		})),
	}

	response, err := client.Do(request)
	Expect(err).NotTo(HaveOccurred())
	defer response.Body.Close()

	return capturedTraceID, capturedTraceParent
}

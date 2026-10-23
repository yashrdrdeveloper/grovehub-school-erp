package response_test

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	oteltrace "go.opentelemetry.io/otel/trace"

	responsedto "grovehub/libs/core/dto/response"
	coreerrors "grovehub/libs/core/errors"
	coretelemetry "grovehub/libs/core/telemetry"
)

const (
	serviceDownCode       = "IDENTITY_SERVICE_DOWN"
	serviceDownMessage    = "Identity service is unavailable"
	invalidRequestCode    = "IDENTITY_INVALID_REQUEST"
	invalidRequestMessage = "Identity request is invalid"
	metadataField         = "email"
)

type metadata struct {
	Field string `json:"field"`
}

var (
	responseTracerProvider *sdktrace.TracerProvider
	responseSpanRecorder   *tracetest.SpanRecorder
)

func TestCoreHTTPResponseUnit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Core HTTP Response Unit Tests")
}

var _ = BeforeEach(func() {
	responseSpanRecorder = tracetest.NewSpanRecorder()
	responseTracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(responseSpanRecorder),
	)
	otel.SetTracerProvider(responseTracerProvider)
})

var _ = AfterEach(func() {
	Expect(responseTracerProvider.Shutdown(context.Background())).To(Succeed())
	otel.SetTracerProvider(oteltrace.NewNoopTracerProvider())
})

func serviceDownDefinition() coreerrors.Definition {
	return coreerrors.NewDefinition(serviceDownCode, serviceDownMessage)
}

func invalidRequestDefinition() coreerrors.Definition {
	return coreerrors.NewDefinition(invalidRequestCode, invalidRequestMessage)
}

func newErrorTraceContext() (context.Context, oteltrace.Span) {
	return coretelemetry.StartSpan(context.Background(), "identity", "test error")
}

func decodeDefaultError(recorder *httptest.ResponseRecorder) responsedto.Error[map[string]any] {
	var payload responsedto.Error[map[string]any]
	Expect(json.Unmarshal(recorder.Body.Bytes(), &payload)).To(Succeed())
	return payload
}

func decodeMetadataError(recorder *httptest.ResponseRecorder) responsedto.Error[metadata] {
	var payload responsedto.Error[metadata]
	Expect(json.Unmarshal(recorder.Body.Bytes(), &payload)).To(Succeed())
	return payload
}

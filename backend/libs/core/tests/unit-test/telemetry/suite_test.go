package telemetry_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	testComponent = "identity"
	testOperation = "health"
)

var testTracerProvider *sdktrace.TracerProvider

func TestCoreTelemetryUnit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Core Telemetry Unit Tests")
}

var _ = BeforeEach(func() {
	testTracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	otel.SetTracerProvider(testTracerProvider)
})

var _ = AfterEach(func() {
	Expect(testTracerProvider.Shutdown(context.Background())).To(Succeed())
	otel.SetTracerProvider(oteltrace.NewNoopTracerProvider())
})

func newJSONLogger() (*slog.Logger, *bytes.Buffer) {
	buffer := &bytes.Buffer{}
	return slog.New(slog.NewJSONHandler(buffer, nil)), buffer
}

func decodeLog(buffer *bytes.Buffer) map[string]any {
	var payload map[string]any
	Expect(json.Unmarshal(buffer.Bytes(), &payload)).To(Succeed())
	return payload
}

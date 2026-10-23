package telemetry

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	defaultServiceName      = "grovehub-platform"
	defaultServiceNamespace = "grovehub"
	defaultEnvironment      = "local"
)

type Shutdown func(context.Context) error

func Initialize(ctx context.Context, logger *slog.Logger) Shutdown {
	serviceName := envOrDefault("OTEL_SERVICE_NAME", defaultServiceName)
	serviceNamespace := envOrDefault("OTEL_SERVICE_NAMESPACE", defaultServiceNamespace)
	environment := envOrDefault("OTEL_DEPLOYMENT_ENVIRONMENT", defaultEnvironment)

	resource := resource.NewSchemaless(
		attribute.String("service.name", serviceName),
		attribute.String("service.namespace", serviceNamespace),
		attribute.String("deployment.environment.name", environment),
	)

	providerOptions := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(resource),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample())),
	}

	exportEnabled := strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")) != ""
	if exportEnabled {
		exporter, err := otlptracehttp.New(ctx)
		if err != nil {
			logger.Error("OpenTelemetry trace exporter initialization failed", "error", err)
			exportEnabled = false
		} else {
			providerOptions = append(providerOptions, sdktrace.WithBatcher(exporter))
		}
	}

	provider := sdktrace.NewTracerProvider(providerOptions...)
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		logger.Error("OpenTelemetry export failed", "error", err)
	}))

	logger.Info(
		"OpenTelemetry initialized",
		"service", serviceName,
		"environment", environment,
		"traceExportEnabled", exportEnabled,
	)

	return provider.Shutdown
}

func envOrDefault(name string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}

	return value
}

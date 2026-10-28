package helpers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	otelruntime "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"grovehub/libs/core/telemetry/contracts"
	"grovehub/libs/core/telemetry/mappers"
	telemetryvalidators "grovehub/libs/core/telemetry/validators"
)

const (
	defaultEnvironment    = "local"
	defaultServiceVersion = "dev"
)

type serviceRuntime struct {
	service             contracts.Service
	tracerProvider      *sdktrace.TracerProvider
	meterProvider       *sdkmetric.MeterProvider
	operationDuration   metric.Float64Histogram
	operationRequests   metric.Int64Counter
	operationErrors     metric.Int64Counter
	httpRequestDuration metric.Float64Histogram
	httpRequestErrors   metric.Int64Counter
	httpActiveRequests  metric.Int64UpDownCounter
}

type Runtime struct {
	services map[string]*serviceRuntime
}

var (
	activeRuntimeMu sync.RWMutex
	activeRuntime   *Runtime
)

func Initialize(ctx context.Context, logger *slog.Logger) contracts.Shutdown {
	environment := envOrDefault("OTEL_DEPLOYMENT_ENVIRONMENT", defaultEnvironment)
	serviceVersion := envOrDefault("OTEL_SERVICE_VERSION", defaultServiceVersion)
	instanceID := envOrDefault("OTEL_SERVICE_INSTANCE_ID", defaultInstanceID())
	traceSampleRatio := configuredTraceSampleRatio(logger)
	metricExportInterval := configuredMetricExportInterval(logger)

	services := []contracts.Service{
		contracts.ServiceLauncherEdge,
		contracts.ServiceFoundation,
		contracts.ServiceIdentity,
	}

	runtime := &Runtime{services: make(map[string]*serviceRuntime, len(services))}
	for _, service := range services {
		serviceRuntime := newServiceRuntime(
			ctx,
			logger,
			service,
			environment,
			serviceVersion,
			instanceID,
			traceSampleRatio,
			metricExportInterval,
		)
		runtime.services[service.Key] = serviceRuntime
	}

	setActiveRuntime(runtime)

	edge := runtime.services[contracts.ServiceLauncherEdge.Key]
	otel.SetTracerProvider(edge.tracerProvider)
	otel.SetMeterProvider(edge.meterProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		logger.Error(
			"OpenTelemetry export failed",
			"service", contracts.ServiceLauncherEdge.Name,
			"error", err,
		)
	}))

	if err := otelruntime.Start(otelruntime.WithMeterProvider(edge.meterProvider)); err != nil {
		logger.Error(
			"OpenTelemetry Go runtime metrics initialization failed",
			"service", contracts.ServiceLauncherEdge.Name,
			"error", err,
		)
	}

	logger.Info(
		"OpenTelemetry initialized",
		"service", contracts.ServiceLauncherEdge.Name,
		"environment", environment,
		"traceSampleRatio", traceSampleRatio,
		"metricExportInterval", metricExportInterval.String(),
	)

	return func(shutdownContext context.Context) error {
		setActiveRuntime(nil)

		var shutdownErrors []error
		for _, service := range services {
			serviceRuntime := runtime.services[service.Key]
			if serviceRuntime == nil {
				continue
			}

			if err := serviceRuntime.meterProvider.Shutdown(shutdownContext); err != nil {
				shutdownErrors = append(
					shutdownErrors,
					fmt.Errorf("shutdown %s meter provider: %w", service.Key, err),
				)
			}
			if err := serviceRuntime.tracerProvider.Shutdown(shutdownContext); err != nil {
				shutdownErrors = append(
					shutdownErrors,
					fmt.Errorf("shutdown %s tracer provider: %w", service.Key, err),
				)
			}
		}
		return errors.Join(shutdownErrors...)
	}
}

func newServiceRuntime(
	ctx context.Context,
	logger *slog.Logger,
	service contracts.Service,
	environment string,
	serviceVersion string,
	instanceID string,
	traceSampleRatio float64,
	metricExportInterval time.Duration,
) *serviceRuntime {
	serviceResource := resource.NewSchemaless(
		attribute.String("service.name", service.Name),
		attribute.String("service.namespace", contracts.ServiceNamespace),
		attribute.String("service.version", serviceVersion),
		attribute.String("service.instance.id", instanceID),
		attribute.String("deployment.environment.name", environment),
		attribute.String("grovehub.service.key", service.Key),
	)

	tracerOptions := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(serviceResource),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(traceSampleRatio))),
	}

	traceExportEnabled := traceExportConfigured()
	if traceExportEnabled {
		exporter, err := otlptracehttp.New(ctx)
		if err != nil {
			logger.Error(
				"OpenTelemetry trace exporter initialization failed",
				"service", service.Name,
				"error", err,
			)
			traceExportEnabled = false
		} else {
			tracerOptions = append(tracerOptions, sdktrace.WithBatcher(exporter))
		}
	}
	tracerProvider := sdktrace.NewTracerProvider(tracerOptions...)

	meterOptions := []sdkmetric.Option{
		sdkmetric.WithResource(serviceResource),
		sdkmetric.WithCardinalityLimit(1000),
	}

	metricExportEnabled := metricExportConfigured()
	if metricExportEnabled {
		exporter, err := otlpmetrichttp.New(ctx)
		if err != nil {
			logger.Error(
				"OpenTelemetry metric exporter initialization failed",
				"service", service.Name,
				"error", err,
			)
			metricExportEnabled = false
		} else {
			meterOptions = append(
				meterOptions,
				sdkmetric.WithReader(
					sdkmetric.NewPeriodicReader(
						exporter,
						sdkmetric.WithInterval(metricExportInterval),
					),
				),
			)
		}
	}
	meterProvider := sdkmetric.NewMeterProvider(meterOptions...)

	serviceRuntime := &serviceRuntime{
		service:        service,
		tracerProvider: tracerProvider,
		meterProvider:  meterProvider,
	}
	serviceRuntime.initializeMetrics(logger)

	logger.Info(
		"OpenTelemetry service initialized",
		"service", service.Name,
		"traceExportEnabled", traceExportEnabled,
		"metricExportEnabled", metricExportEnabled,
	)

	return serviceRuntime
}

func TracerProvider(serviceName string) trace.TracerProvider {
	service := mappers.ResolveService(serviceName)
	if runtime := getActiveRuntime(); runtime != nil {
		if serviceRuntime := runtime.services[service.Key]; serviceRuntime != nil {
			return serviceRuntime.tracerProvider
		}
	}
	return otel.GetTracerProvider()
}

func MeterProvider(serviceName string) metric.MeterProvider {
	service := mappers.ResolveService(serviceName)
	if runtime := getActiveRuntime(); runtime != nil {
		if serviceRuntime := runtime.services[service.Key]; serviceRuntime != nil {
			return serviceRuntime.meterProvider
		}
	}
	return otel.GetMeterProvider()
}

func setActiveRuntime(runtime *Runtime) {
	activeRuntimeMu.Lock()
	defer activeRuntimeMu.Unlock()
	activeRuntime = runtime
}

func getActiveRuntime() *Runtime {
	activeRuntimeMu.RLock()
	defer activeRuntimeMu.RUnlock()
	return activeRuntime
}

func configuredTraceSampleRatio(logger *slog.Logger) float64 {
	value := strings.TrimSpace(os.Getenv("GROVEHUB_TRACE_SAMPLE_RATIO"))
	ratio, err := telemetryvalidators.TraceSampleRatio(value)
	if err == nil {
		return ratio
	}
	logger.Warn(
		"invalid GROVEHUB_TRACE_SAMPLE_RATIO; using default",
		"value", value,
		"default", telemetryvalidators.DefaultTraceSampleRatio,
	)
	return telemetryvalidators.DefaultTraceSampleRatio
}

func configuredMetricExportInterval(logger *slog.Logger) time.Duration {
	value := strings.TrimSpace(os.Getenv("GROVEHUB_METRIC_EXPORT_INTERVAL"))
	interval, err := telemetryvalidators.MetricExportInterval(value)
	if err == nil {
		return interval
	}
	logger.Warn(
		"invalid GROVEHUB_METRIC_EXPORT_INTERVAL; using default",
		"value", value,
		"default", telemetryvalidators.DefaultMetricExportInterval.String(),
	)
	return telemetryvalidators.DefaultMetricExportInterval
}

func traceExportConfigured() bool {
	return firstNonEmptyEnvironmentValue(
		"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT",
		"OTEL_EXPORTER_OTLP_ENDPOINT",
	) != ""
}

func metricExportConfigured() bool {
	return firstNonEmptyEnvironmentValue(
		"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT",
		"OTEL_EXPORTER_OTLP_ENDPOINT",
	) != ""
}

func firstNonEmptyEnvironmentValue(names ...string) string {
	for _, name := range names {
		if value := strings.TrimSpace(os.Getenv(name)); value != "" {
			return value
		}
	}
	return ""
}

func defaultInstanceID() string {
	hostname, err := os.Hostname()
	if err != nil || strings.TrimSpace(hostname) == "" {
		hostname = "unknown-host"
	}
	return fmt.Sprintf("%s-%d", hostname, os.Getpid())
}

func envOrDefault(name string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	return value
}

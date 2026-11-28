package runtime

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
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

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/telemetry/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/telemetry/internal/src/helpers"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/telemetry/internal/src/config"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/telemetry/internal/src/constants"
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
	environment := config.EnvOrDefault("OTEL_DEPLOYMENT_ENVIRONMENT", constants.DefaultEnvironment)
	serviceVersion := config.EnvOrDefault("OTEL_SERVICE_VERSION", constants.DefaultServiceVersion)
	instanceID := config.EnvOrDefault("OTEL_SERVICE_INSTANCE_ID", config.DefaultInstanceID())
	traceSampleRatio := config.ConfiguredTraceSampleRatio(logger)
	metricExportInterval := config.ConfiguredMetricExportInterval(logger)

	services := []contracts.Service{
		contracts.ServiceLauncherEdge,
		contracts.ServiceFoundation,
		contracts.ServiceIdentity,
		contracts.ServiceAudit,
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
		attribute.String("service.namespace", constants.ServiceNamespace),
		attribute.String("service.version", serviceVersion),
		attribute.String("service.instance.id", instanceID),
		attribute.String("deployment.environment.name", environment),
		attribute.String("grovehub.service.key", service.Key),
	)

	tracerOptions := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(serviceResource),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(traceSampleRatio))),
	}

	traceExportEnabled := config.TraceExportConfigured()
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
		sdkmetric.WithCardinalityLimit(constants.MetricCardinalityLimit),
	}

	metricExportEnabled := config.MetricExportConfigured()
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
	service := helpers.ResolveService(serviceName)
	if runtime := getActiveRuntime(); runtime != nil {
		if serviceRuntime := runtime.services[service.Key]; serviceRuntime != nil {
			return serviceRuntime.tracerProvider
		}
	}
	return otel.GetTracerProvider()
}

func MeterProvider(serviceName string) metric.MeterProvider {
	service := helpers.ResolveService(serviceName)
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


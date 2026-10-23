package telemetry

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var latencyBoundariesSeconds = []float64{
	0.005,
	0.01,
	0.025,
	0.05,
	0.075,
	0.1,
	0.25,
	0.5,
	0.75,
	1,
	2.5,
	5,
	7.5,
	10,
}

func (runtime *serviceRuntime) initializeMetrics(logger *slog.Logger) {
	meter := runtime.meterProvider.Meter("grovehub/" + runtime.service.Key)

	operationDuration, err := meter.Float64Histogram(
		"grovehub.service.operation.duration",
		metric.WithDescription("Duration of a GroveHub service operation."),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(latencyBoundariesSeconds...),
	)
	if err != nil {
		logger.Error(
			"failed to create service operation duration metric",
			"service", runtime.service.Name,
			"error", err,
		)
	} else {
		runtime.operationDuration = operationDuration
	}

	operationRequests, err := meter.Int64Counter(
		"grovehub.service.operation.requests",
		metric.WithDescription("Number of GroveHub service operations."),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		logger.Error(
			"failed to create service operation request metric",
			"service", runtime.service.Name,
			"error", err,
		)
	} else {
		runtime.operationRequests = operationRequests
	}

	operationErrors, err := meter.Int64Counter(
		"grovehub.service.operation.errors",
		metric.WithDescription("Number of failed GroveHub service operations."),
		metric.WithUnit("{error}"),
	)
	if err != nil {
		logger.Error(
			"failed to create service operation error metric",
			"service", runtime.service.Name,
			"error", err,
		)
	} else {
		runtime.operationErrors = operationErrors
	}

	httpRequestDuration, err := meter.Float64Histogram(
		"grovehub.http.server.request.duration",
		metric.WithDescription("Duration of GroveHub HTTP server requests by matched route."),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(latencyBoundariesSeconds...),
	)
	if err != nil {
		logger.Error(
			"failed to create route-aware HTTP request duration metric",
			"service", runtime.service.Name,
			"error", err,
		)
	} else {
		runtime.httpRequestDuration = httpRequestDuration
	}

	httpRequestErrors, err := meter.Int64Counter(
		"grovehub.http.server.request.errors",
		metric.WithDescription("Number of failed GroveHub HTTP server requests by matched route."),
		metric.WithUnit("{error}"),
	)
	if err != nil {
		logger.Error(
			"failed to create route-aware HTTP request error metric",
			"service", runtime.service.Name,
			"error", err,
		)
	} else {
		runtime.httpRequestErrors = httpRequestErrors
	}

	httpActiveRequests, err := meter.Int64UpDownCounter(
		"grovehub.http.server.active_requests",
		metric.WithDescription("Number of GroveHub HTTP server requests currently in flight."),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		logger.Error(
			"failed to create active HTTP request metric",
			"service", runtime.service.Name,
			"error", err,
		)
	} else {
		runtime.httpActiveRequests = httpActiveRequests
	}

	_, err = meter.Int64ObservableGauge(
		"grovehub.service.up",
		metric.WithDescription("Reports 1 while the GroveHub logical service is running."),
		metric.WithUnit("1"),
		metric.WithInt64Callback(func(_ context.Context, observer metric.Int64Observer) error {
			observer.Observe(1)
			return nil
		}),
	)
	if err != nil {
		logger.Error(
			"failed to create service availability metric",
			"service", runtime.service.Name,
			"error", err,
		)
	}
}

func RecordOperation(
	ctx context.Context,
	serviceName string,
	operation string,
	duration time.Duration,
	statusCode int,
) {
	service := ResolveService(serviceName)
	runtime := getActiveRuntime()
	if runtime == nil {
		return
	}

	serviceRuntime := runtime.services[service.Key]
	if serviceRuntime == nil {
		return
	}

	result := requestResult(statusCode)
	attributes := []attribute.KeyValue{
		attribute.String("grovehub.operation", operation),
		attribute.String("grovehub.result", result),
		attribute.Int("http.response.status_code", statusCode),
	}

	if serviceRuntime.operationDuration != nil {
		serviceRuntime.operationDuration.Record(
			ctx,
			duration.Seconds(),
			metric.WithAttributes(attributes...),
		)
	}
	if serviceRuntime.operationRequests != nil {
		serviceRuntime.operationRequests.Add(
			ctx,
			1,
			metric.WithAttributes(attributes...),
		)
	}
	if result == "error" && serviceRuntime.operationErrors != nil {
		serviceRuntime.operationErrors.Add(
			ctx,
			1,
			metric.WithAttributes(attributes...),
		)
	}
}

func RecordHTTPServerRequest(
	ctx context.Context,
	serviceName string,
	method string,
	route string,
	duration time.Duration,
	statusCode int,
) {
	service := ResolveService(serviceName)
	runtime := getActiveRuntime()
	if runtime == nil {
		return
	}

	serviceRuntime := runtime.services[service.Key]
	if serviceRuntime == nil {
		return
	}

	result := requestResult(statusCode)
	attributes := []attribute.KeyValue{
		attribute.String("http.request.method", method),
		attribute.String("http.route", route),
		attribute.Int("http.response.status_code", statusCode),
		attribute.String("grovehub.result", result),
	}

	if serviceRuntime.httpRequestDuration != nil {
		serviceRuntime.httpRequestDuration.Record(
			ctx,
			duration.Seconds(),
			metric.WithAttributes(attributes...),
		)
	}
	if result == "error" && serviceRuntime.httpRequestErrors != nil {
		serviceRuntime.httpRequestErrors.Add(
			ctx,
			1,
			metric.WithAttributes(attributes...),
		)
	}
}

func AddActiveHTTPServerRequest(
	ctx context.Context,
	serviceName string,
	method string,
	delta int64,
) {
	service := ResolveService(serviceName)
	runtime := getActiveRuntime()
	if runtime == nil {
		return
	}

	serviceRuntime := runtime.services[service.Key]
	if serviceRuntime == nil || serviceRuntime.httpActiveRequests == nil {
		return
	}

	serviceRuntime.httpActiveRequests.Add(
		ctx,
		delta,
		metric.WithAttributes(attribute.String("http.request.method", method)),
	)
}

func requestResult(statusCode int) string {
	if statusCode >= http.StatusBadRequest {
		return "error"
	}
	return "success"
}

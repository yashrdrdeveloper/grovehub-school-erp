package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"grovehub/libs/core/telemetry/internal/src/constants"
)

func ConfiguredTraceSampleRatio(logger *slog.Logger) float64 {
	value := strings.TrimSpace(os.Getenv("GROVEHUB_TRACE_SAMPLE_RATIO"))
	ratio, err := TraceSampleRatio(value)
	if err == nil {
		return ratio
	}
	logger.Warn(
		"invalid GROVEHUB_TRACE_SAMPLE_RATIO; using default",
		"value", value,
		"default", constants.DefaultTraceSampleRatio,
	)
	return constants.DefaultTraceSampleRatio
}

func ConfiguredMetricExportInterval(logger *slog.Logger) time.Duration {
	value := strings.TrimSpace(os.Getenv("GROVEHUB_METRIC_EXPORT_INTERVAL"))
	interval, err := MetricExportInterval(value)
	if err == nil {
		return interval
	}
	logger.Warn(
		"invalid GROVEHUB_METRIC_EXPORT_INTERVAL; using default",
		"value", value,
		"default", constants.DefaultMetricExportInterval.String(),
	)
	return constants.DefaultMetricExportInterval
}

func TraceExportConfigured() bool {
	return firstNonEmptyEnvironmentValue(
		"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT",
		"OTEL_EXPORTER_OTLP_ENDPOINT",
	) != ""
}

func MetricExportConfigured() bool {
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

func DefaultInstanceID() string {
	return InstanceID(os.Hostname, os.Getpid())
}

// InstanceID resolves a process identity using the platform hostname boundary.
func InstanceID(host func() (string, error), pid int) string {
	hostname, err := host()
	if err != nil || strings.TrimSpace(hostname) == "" {
		hostname = "unknown-host"
	}
	return fmt.Sprintf("%s-%d", hostname, pid)
}

func EnvOrDefault(name string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	return value
}

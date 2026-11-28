package constants

import "time"

const (
	ServiceNamespace = "grovehub"
	DefaultEnvironment = "local"
	DefaultServiceVersion = "dev"
	DefaultTraceSampleRatio = 1.0
	DefaultMetricExportInterval = 15 * time.Second
	MinMetricExportInterval = 5 * time.Second
	MetricCardinalityLimit = 1000
)

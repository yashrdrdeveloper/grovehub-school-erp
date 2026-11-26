package validators

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultTraceSampleRatio = 1.0
)

const DefaultMetricExportInterval = 15 * time.Second

func TraceSampleRatio(value string) (float64, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return DefaultTraceSampleRatio, nil
	}

	ratio, err := strconv.ParseFloat(value, 64)
	if err != nil || ratio < 0 || ratio > 1 {
		return DefaultTraceSampleRatio, fmt.Errorf("trace sample ratio must be between 0 and 1")
	}
	return ratio, nil
}

func MetricExportInterval(value string) (time.Duration, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return DefaultMetricExportInterval, nil
	}

	interval, err := time.ParseDuration(value)
	if err != nil || interval < 5*time.Second {
		return DefaultMetricExportInterval, fmt.Errorf("metric export interval must be at least 5s")
	}
	return interval, nil
}

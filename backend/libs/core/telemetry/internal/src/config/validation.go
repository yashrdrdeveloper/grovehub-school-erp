package config

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/telemetry/internal/src/constants"
)

func TraceSampleRatio(value string) (float64, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return constants.DefaultTraceSampleRatio, nil
	}

	ratio, err := strconv.ParseFloat(value, 64)
	if err != nil || ratio < 0 || ratio > 1 {
		return constants.DefaultTraceSampleRatio, fmt.Errorf("trace sample ratio must be between 0 and 1")
	}
	return ratio, nil
}

func MetricExportInterval(value string) (time.Duration, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return constants.DefaultMetricExportInterval, nil
	}

	interval, err := time.ParseDuration(value)
	if err != nil || interval < constants.MinMetricExportInterval {
		return constants.DefaultMetricExportInterval, fmt.Errorf("metric export interval must be at least 5s")
	}
	return interval, nil
}

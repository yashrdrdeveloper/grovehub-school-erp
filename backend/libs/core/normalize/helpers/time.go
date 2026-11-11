package helpers

import (
	"errors"
	"strings"
	"time"
)

var ErrInvalidTime = errors.New("invalid RFC3339 time")

// RFC3339UTC normalizes a timestamp to UTC with nanosecond precision.
func RFC3339UTC(value string) (string, error) {
	parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(value))
	if err != nil {
		return "", ErrInvalidTime
	}
	return parsed.UTC().Format(time.RFC3339Nano), nil
}

package normalize

import (
	"strings"
	"time"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/utils/internal/src/contracts"
)

// NormalizeRFC3339UTC normalizes a timestamp to UTC with nanosecond precision.
func NormalizeRFC3339UTC(value string) (string, error) {
	parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(value))
	if err != nil {
		return "", contracts.ErrInvalidTime
	}
	return parsed.UTC().Format(time.RFC3339Nano), nil
}

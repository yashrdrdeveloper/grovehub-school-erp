package normalize

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"grovehub/libs/core/utils/internal/src/contracts"
)

// ValidateUUID requires a nonzero UUID in canonical lowercase form, without normalization.
func ValidateUUID(value string) error {
	id, err := uuid.Parse(value)
	if err != nil || id == uuid.Nil || id.String() != value {
		return contracts.ErrInvalidUUID
	}
	return nil
}

// NormalizeText trims surrounding whitespace and enforces a nonempty byte limit.
func NormalizeText(value string, maxBytes int) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" || len(value) > maxBytes {
		return "", contracts.ErrInvalidText
	}
	return value, nil
}

// NormalizeStrictRFC3339UTC preserves the strict no-whitespace timestamp contract.
func NormalizeStrictRFC3339UTC(value string) (string, error) {
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return "", contracts.ErrInvalidTime
	}
	return parsed.UTC().Format(time.RFC3339Nano), nil
}

package key

import (
	"strings"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/idempotency/internal/src/constants"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/idempotency/internal/src/contracts"
)

// ValidateKey trims surrounding whitespace and validates the existing ASCII allowlist.
func ValidateKey(value string) (string, error) {
	key := strings.TrimSpace(value)
	if key == "" || len(key) > constants.MaxKeyLength {
		return "", contracts.ErrInvalidKey
	}
	for _, character := range key {
		letter := character >= 'A' && character <= 'Z' || character >= 'a' && character <= 'z'
		digit := character >= '0' && character <= '9'
		if letter || digit {
			continue
		}
		switch character {
		case '-', '_', '.', ':':
			continue
		default:
			return "", contracts.ErrInvalidKey
		}
	}
	return key, nil
}

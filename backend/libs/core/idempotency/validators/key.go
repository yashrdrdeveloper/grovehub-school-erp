package validators

import (
	"strings"

	"grovehub/libs/core/idempotency/contracts"
)

func ValidateKey(value string) (string, error) {
	key := strings.TrimSpace(value)
	if key == "" || len(key) > contracts.MaxKeyLength {
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

package normalize

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

var ErrInvalidUUIDv7 = errors.New("invalid UUIDv7")

func UUIDv7(value string) (string, error) {
	normalized := strings.TrimSpace(value)
	parsed, err := uuid.Parse(normalized)
	if err != nil || parsed == uuid.Nil || parsed.Version() != uuid.Version(7) {
		return "", ErrInvalidUUIDv7
	}
	return parsed.String(), nil
}

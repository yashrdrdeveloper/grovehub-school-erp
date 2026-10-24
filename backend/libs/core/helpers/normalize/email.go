package normalize

import (
	"errors"
	"net/mail"
	"strings"
)

const maxEmailLength = 254

var ErrInvalidEmail = errors.New("invalid email")

func Email(value string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" || len(normalized) > maxEmailLength {
		return "", ErrInvalidEmail
	}

	address, err := mail.ParseAddress(normalized)
	if err != nil || strings.ToLower(address.Address) != normalized {
		return "", ErrInvalidEmail
	}
	return normalized, nil
}

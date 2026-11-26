package helpers

import (
	"errors"
	"net/mail"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

const maxEmailLength = 254

var (
	ErrInvalidEmail  = errors.New("invalid email")
	ErrInvalidPhone  = errors.New("invalid phone")
	ErrInvalidUUIDv7 = errors.New("invalid UUIDv7")
	phonePattern     = regexp.MustCompile(`^\+[1-9][0-9]{7,14}$`)
)

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

func Phone(value string) (string, error) {
	normalized, err := PhoneCharacters(value)
	if err != nil || !phonePattern.MatchString(normalized) {
		return "", ErrInvalidPhone
	}
	return normalized, nil
}

func PhoneCharacters(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	builder := strings.Builder{}

	for index, character := range trimmed {
		digit := character >= '0' && character <= '9'
		separator := character == ' ' || character == '-' || character == '(' || character == ')'

		if character == '+' && index == 0 {
			builder.WriteRune(character)
		} else if digit {
			builder.WriteRune(character)
		} else if !separator {
			return "", ErrInvalidPhone
		}
	}

	return builder.String(), nil
}

func UUIDv7(value string) (string, error) {
	normalized := strings.TrimSpace(value)
	parsed, err := uuid.Parse(normalized)
	if err != nil || parsed == uuid.Nil || parsed.Version() != uuid.Version(7) {
		return "", ErrInvalidUUIDv7
	}
	return parsed.String(), nil
}

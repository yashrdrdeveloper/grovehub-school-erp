package normalize

import (
	"net/mail"
	"regexp"
	"strings"

	"github.com/google/uuid"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/utils/internal/src/constants"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/utils/internal/src/contracts"
)

var phonePattern = regexp.MustCompile(constants.PhonePattern)

func NormalizeEmail(value string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" || len(normalized) > constants.MaxEmailLength {
		return "", contracts.ErrInvalidEmail
	}

	address, err := mail.ParseAddress(normalized)
	if err != nil || strings.ToLower(address.Address) != normalized {
		return "", contracts.ErrInvalidEmail
	}
	return normalized, nil
}

func NormalizePhone(value string) (string, error) {
	normalized, err := NormalizePhoneCharacters(value)
	if err != nil || !phonePattern.MatchString(normalized) {
		return "", contracts.ErrInvalidPhone
	}
	return normalized, nil
}

func NormalizePhoneCharacters(value string) (string, error) {
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
			return "", contracts.ErrInvalidPhone
		}
	}

	return builder.String(), nil
}

func NormalizeUUIDv7(value string) (string, error) {
	normalized := strings.TrimSpace(value)
	parsed, err := uuid.Parse(normalized)
	if err != nil || parsed == uuid.Nil || parsed.Version() != uuid.Version(7) {
		return "", contracts.ErrInvalidUUIDv7
	}
	return parsed.String(), nil
}

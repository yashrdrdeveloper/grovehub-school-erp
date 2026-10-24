package normalize

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrInvalidPhone = errors.New("invalid phone")
	phonePattern    = regexp.MustCompile(`^\+[1-9][0-9]{7,14}$`)
)

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

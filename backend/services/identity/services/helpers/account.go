package helpers

import (
	"net/mail"
	"regexp"
	"strings"

	"github.com/google/uuid"

	coreerrors "grovehub/libs/core/errors"
	identityerror "grovehub/services/identity/constants/error"
	"grovehub/services/identity/entity"
)

const maxEmailLength = 254

var phonePattern = regexp.MustCompile(`^\+[1-9][0-9]{7,14}$`)

func NormalizeAccountId(value string) (entity.AccountId, error) {
	normalized := strings.TrimSpace(value)
	parsed, err := uuid.Parse(normalized)
	if err != nil || parsed == uuid.Nil || parsed.Version() != uuid.Version(7) {
		return "", invalidAccount()
	}
	return entity.AccountId(parsed.String()), nil
}

func NormalizePersonGid(value string) (entity.PersonGid, error) {
	normalized := strings.TrimSpace(value)
	parsed, err := uuid.Parse(normalized)
	if err != nil || parsed == uuid.Nil || parsed.Version() != uuid.Version(7) {
		return "", invalidAccount()
	}
	return entity.PersonGid(parsed.String()), nil
}

func NormalizeEmail(value string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" || len(normalized) > maxEmailLength {
		return "", invalidAccount()
	}
	address, err := mail.ParseAddress(normalized)
	if err != nil || strings.ToLower(address.Address) != normalized {
		return "", invalidAccount()
	}
	return normalized, nil
}

func NormalizePhone(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", nil
	}

	builder := strings.Builder{}
	for index, character := range trimmed {
		digit := character >= '0' && character <= '9'
		separator := character == ' ' || character == '-' || character == '(' || character == ')'
		if character == '+' && index == 0 {
			builder.WriteRune(character)
		} else if digit {
			builder.WriteRune(character)
		} else if !separator {
			return "", invalidAccount()
		}
	}

	normalized := builder.String()
	if !phonePattern.MatchString(normalized) {
		return "", invalidAccount()
	}
	return normalized, nil
}

func invalidAccount() error {
	return coreerrors.NewError(identityerror.InvalidAccount)
}

package helpers

import (
	"unicode/utf8"

	coreerrors "grovehub/libs/core/errors"
	identityerror "grovehub/services/identity/constants/error"
)

const (
	minPasswordLength = 8
	maxPasswordLength = 128
)

func ValidatePassword(value string) error {
	if err := validatePasswordCandidate(value); err != nil {
		return err
	}
	if utf8.RuneCountInString(value) < minPasswordLength {
		return invalidPassword()
	}
	return nil
}

func ValidatePasswordCandidate(value string) error {
	return validatePasswordCandidate(value)
}

func validatePasswordCandidate(value string) error {
	if !utf8.ValidString(value) {
		return invalidPassword()
	}
	length := utf8.RuneCountInString(value)
	if length < 1 || length > maxPasswordLength {
		return invalidPassword()
	}
	return nil
}

func invalidPassword() error {
	return coreerrors.NewError(identityerror.InvalidPassword)
}

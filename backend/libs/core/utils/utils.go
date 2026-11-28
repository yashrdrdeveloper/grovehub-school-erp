// Package utils exposes domain-agnostic conversions and normalization.
package utils

import (
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/utils/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/utils/internal/src/operations/normalize"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/utils/internal/src/operations/strings"
)

// OptionalString converts an empty string to nil for nullable fields.
func OptionalString(value string) *string {
	return strings.OptionalString(value)
}

// NullableString dereferences a nullable string, defaulting nil to "".
func NullableString(value *string) string {
	return strings.NullableString(value)
}

var (
	ErrInvalidEmail  = contracts.ErrInvalidEmail
	ErrInvalidPhone  = contracts.ErrInvalidPhone
	ErrInvalidUUIDv7 = contracts.ErrInvalidUUIDv7
	ErrInvalidUUID   = contracts.ErrInvalidUUID
	ErrInvalidText   = contracts.ErrInvalidText
	ErrInvalidTime   = contracts.ErrInvalidTime
)

func NormalizeEmail(value string) (string, error) {
	return normalize.NormalizeEmail(value)
}

func NormalizePhone(value string) (string, error) {
	return normalize.NormalizePhone(value)
}

func NormalizePhoneCharacters(value string) (string, error) {
	return normalize.NormalizePhoneCharacters(value)
}

func NormalizeUUIDv7(value string) (string, error) {
	return normalize.NormalizeUUIDv7(value)
}

func NormalizeRFC3339UTC(value string) (string, error) {
	return normalize.NormalizeRFC3339UTC(value)
}

func ValidateUUID(value string) error {
	return normalize.ValidateUUID(value)
}

func NormalizeText(value string, maxBytes int) (string, error) {
	return normalize.NormalizeText(value, maxBytes)
}

func NormalizeStrictRFC3339UTC(value string) (string, error) {
	return normalize.NormalizeStrictRFC3339UTC(value)
}

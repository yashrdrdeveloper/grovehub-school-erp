// Package utils exposes small, domain-agnostic conversions shared by services.
package utils

import "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/utils/internal/src/operations/strings"

// OptionalString converts an empty string to nil for nullable fields.
func OptionalString(value string) *string {
	return strings.OptionalString(value)
}

// NullableString dereferences a nullable string, defaulting nil to "".
func NullableString(value *string) string {
	return strings.NullableString(value)
}

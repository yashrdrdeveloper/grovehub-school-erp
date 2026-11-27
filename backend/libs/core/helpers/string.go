// Package helpers holds small, domain-agnostic conversions shared by services.
package helpers

// OptionalString converts an empty string to nil for nullable fields.
func OptionalString(value string) *string {
    if value == "" { return nil }
    return &value
}

// NullableString dereferences a nullable string, defaulting nil to "".
func NullableString(value *string) string {
    if value == nil { return "" }
    return *value
}

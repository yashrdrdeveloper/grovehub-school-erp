package strings

// OptionalString converts an empty string to nil without trimming or normalizing.
func OptionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

// NullableString returns the value of a nonnil pointer or the empty string.
func NullableString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

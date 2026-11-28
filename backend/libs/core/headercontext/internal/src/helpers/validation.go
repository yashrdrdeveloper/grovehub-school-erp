package helpers

import "grovehub/libs/core/headercontext/internal/src/constants"

// ValidID accepts an absent identifier or a bounded ASCII metadata identifier.
func ValidID(value string) bool {
	if len(value) > constants.MaxMetadataIDLength {
		return false
	}
	for _, c := range value {
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' {
			continue
		}
		switch c {
		case '-', '_', '.', ':', '/':
			continue
		default:
			return false
		}
	}
	return true
}

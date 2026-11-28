package helpers

import (
	"strings"

	"grovehub/services/registry/internal/app/services/authorization/constants"
	"grovehub/services/registry/internal/app/validators/workload"
)

// KeyID reads only an untrusted lookup hint; it authenticates nothing.
func KeyID(token string) (string, bool) {
	if len(token) == 0 || len(token) > constants.MaxTokenBytes {
		return "", false
	}
	parts := strings.Split(token, ".")
	if len(parts) != constants.TokenParts || parts[0] != constants.Version ||
		workload.ValidateKeyID(parts[1]) != nil {
		return "", false
	}
	return parts[1], true
}

package helpers

import (
	"net/http"
	"strings"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/headercontext/internal/src/contracts"
)

// ReadOne reads an optional single, nonblank header and trims surrounding space.
func ReadOne(header http.Header, name string) (string, error) {
	values := header.Values(name)
	if len(values) == 0 {
		return "", nil
	}
	if len(values) != 1 {
		return "", contracts.ErrInvalidHeader
	}
	value := strings.TrimSpace(values[0])
	if value == "" {
		return "", contracts.ErrInvalidHeader
	}
	return value, nil
}

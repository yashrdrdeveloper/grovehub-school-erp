package definitions

import (
	"net/http"

	coreerrors "grovehub/libs/core/errors"
)

var (
	Forbidden   = coreerrors.NewDefinition("REGISTRY_FORBIDDEN", "Access denied", http.StatusForbidden)
	Invalid     = coreerrors.NewDefinition("REGISTRY_INVALID", "Invalid registry request", http.StatusBadRequest)
	NotFound    = coreerrors.NewDefinition("REGISTRY_NOT_FOUND", "Registry resource not found", http.StatusNotFound)
	Conflict    = coreerrors.NewDefinition("REGISTRY_CONFLICT", "Registry resource or version conflict", http.StatusConflict)
	Unavailable = coreerrors.NewDefinition("REGISTRY_UNAVAILABLE", "Registry unavailable", http.StatusServiceUnavailable)
)

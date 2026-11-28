package helpers

import (
	"net/http"

	coreerrors "grovehub/libs/core/errors"
)

var (
	unauthorized = coreerrors.NewDefinition("GATEWAY_UNAUTHENTICATED", "Authentication required", http.StatusUnauthorized)
	forbidden    = coreerrors.NewDefinition("GATEWAY_FORBIDDEN", "Access denied", http.StatusForbidden)
	rateLimited  = coreerrors.NewDefinition("GATEWAY_RATE_LIMITED", "Rate limit exceeded", http.StatusTooManyRequests)
)

func gatewayDefinition(status int) coreerrors.Definition {
	switch status {
	case http.StatusUnauthorized:
		return unauthorized
	case http.StatusTooManyRequests:
		return rateLimited
	default:
		return forbidden
	}
}

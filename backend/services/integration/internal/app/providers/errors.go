package providers

import (
	"errors"
	"net/http"
)

// Error exposes a safe normalized provider outcome.
// Never include raw upstream bodies, authorization headers, or tokens.
type Error struct {
	Code       string
	HTTPStatus int
	Retryable  bool
}

func (e *Error) Error() string { return "provider operation failed: " + e.Code }

func NormalizeStatus(status int) error {
	switch status {
	case http.StatusBadRequest:
		return &Error{Code: "INVALID_REQUEST", HTTPStatus: http.StatusBadRequest}
	case http.StatusUnauthorized:
		return &Error{Code: "UNAUTHORIZED", HTTPStatus: http.StatusUnauthorized}
	case http.StatusForbidden:
		return &Error{Code: "FORBIDDEN", HTTPStatus: http.StatusForbidden}
	case http.StatusNotFound:
		return &Error{Code: "NOT_FOUND", HTTPStatus: http.StatusNotFound}
	case http.StatusConflict:
		return &Error{Code: "CONFLICT", HTTPStatus: http.StatusConflict}
	case http.StatusTooManyRequests:
		return &Error{Code: "RATE_LIMITED", HTTPStatus: http.StatusTooManyRequests, Retryable: true}
	default:
		if status >= 500 {
			return &Error{Code: "UPSTREAM_UNAVAILABLE", HTTPStatus: http.StatusBadGateway, Retryable: true}
		}
		return &Error{Code: "UPSTREAM_REJECTED", HTTPStatus: http.StatusBadGateway}
	}
}
func HTTPStatus(err error) int {
	var e *Error
	if errors.As(err, &e) {
		return e.HTTPStatus
	}
	return 0
}

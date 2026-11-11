// Package rest decodes, normalizes and validates typed HTTP request bodies.
// Service execution and business rules remain the controller/service's responsibility.
package rest

import (
	"fmt"
	"net/http"

	httprequest "grovehub/libs/core/http/request"
)

// Config defines the ordered REST input stages. Optional callbacks are skipped.
// Decode defaults to the existing permissive JSON decoder; endpoints requiring
// additionalProperties=false should supply httprequest.DecodeStrictJSON.
type Config[T any] struct {
	Decode    func(*http.Request, any) error
	Normalize func(T) (T, error)
	Validate  func(T) error
}

// DecodeError identifies malformed HTTP input without exposing its contents.
// Domain validators return their existing errors unchanged.
type DecodeError struct{ Err error }

func (err *DecodeError) Error() string { return fmt.Sprintf("decode request: %v", err.Err) }
func (err *DecodeError) Unwrap() error { return err.Err }

// Run executes decode -> normalize -> validate, failing at the first error.
// No service invocation, HTTP response writing or request-context mutation occurs.
func Run[T any](request *http.Request, config Config[T]) (T, error) {
	var value T
	decode := config.Decode
	if decode == nil {
		decode = httprequest.DecodeJSON
	}
	if err := decode(request, &value); err != nil {
		return value, &DecodeError{Err: err}
	}
	if config.Normalize != nil {
		normalized, err := config.Normalize(value)
		if err != nil {
			var zero T
			return zero, err
		}
		value = normalized
	}
	if config.Validate != nil {
		if err := config.Validate(value); err != nil {
			var zero T
			return zero, err
		}
	}
	return value, nil
}

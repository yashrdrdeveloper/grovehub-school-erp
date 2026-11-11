// Package rest decodes HTTP JSON and executes the shared ordered DTO
// normalization/validation pipeline before the controller calls a service.
package rest

import (
	"fmt"
	"net/http"

	annotationvalidation "grovehub/libs/core/annotation/validation"
	httprequest "grovehub/libs/core/http/request"
)

// Config controls transport decoding and resource-owned validation rules.
// Request DTO normalization and tag validation always run automatically.
// A route may opt into strict JSON with DecodeStrictJSON.
type Config[T any] struct {
	Decode func(*http.Request, any) error
	Rules  annotationvalidation.Rules[T]
}

// DecodeError distinguishes malformed HTTP JSON from validation failures.
type DecodeError struct{ Err error }

func (err *DecodeError) Error() string { return fmt.Sprintf("decode request: %v", err.Err) }
func (err *DecodeError) Unwrap() error { return err.Err }

// Run follows decode -> tag normalization -> resource normalization ->
// DTO tag validation -> custom resource validation. It returns the validated
// DTO to its controller; it never invokes the service or writes a response.
func Run[T any](request *http.Request, config Config[T]) (T, error) {
	var value T
	decode := config.Decode
	if decode == nil {
		decode = httprequest.DecodeJSON
	}
	if err := decode(request, &value); err != nil {
		return value, &DecodeError{Err: err}
	}
	return annotationvalidation.Process(value, config.Rules)
}

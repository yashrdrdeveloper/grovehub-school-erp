// Package rest binds and validates typed HTTP inputs and executes service
// callbacks. Controllers own composition; services never depend on net/http.
package rest

import (
	"context"
	"fmt"
	"net/http"

	validation "grovehub/libs/core/annotation/validation"
	coreerrors "grovehub/libs/core/errors"
	httprequest "grovehub/libs/core/http/request"
)

// Empty marks an absent path, query, or JSON body.
type Empty struct{}

// Input contains the normalized and validated values passed to a service.
type Input[Path, Query, Body any] struct {
	Path  Path
	Query Query
	Body  Body
}

// Result is an HTTP presentation of a successful service result.
type Result[T any] struct {
	StatusCode int
	Data       T
	Headers    http.Header
}

func OK[T any](data T) Result[T] { return Result[T]{StatusCode: http.StatusOK, Data: data} }
func Created[T any](data T) Result[T] { return Result[T]{StatusCode: http.StatusCreated, Data: data} }
func CreatedOrOK[T any](data T, created bool) Result[T] {
	if created { return Created(data) }
	return OK(data)
}

// Config binds independent transport sources and runs all validation before
// calling Service. Invalid* definitions are for transport/DTO input failures;
// Fallback is used only when a service or pipeline error is not classified.
type Config[Path, Query, Body, Response any] struct {
	Path        validation.Rules[Path]
	Query       validation.Rules[Query]
	Body        validation.Rules[Body]
	Validate    func(Input[Path, Query, Body]) error
	Service     func(context.Context, Input[Path, Query, Body]) (Result[Response], error)
	Decode      func(*http.Request, any) error
	InvalidPath coreerrors.Definition
	InvalidQuery coreerrors.Definition
	InvalidBody coreerrors.Definition
	Fallback    coreerrors.Definition
	Writer      Writer[Response]
}

// DecodeError reports a bad request body without reflecting user-supplied data.
type DecodeError struct{ Err error }

func (err *DecodeError) Error() string { return fmt.Sprintf("decode request: %v", err.Err) }
func (err *DecodeError) Unwrap() error { return err.Err }

// ReadBody remains available for callers needing only decoding and validation.
// Prefer Run for HTTP controllers, so error handling and response output are shared.
func ReadBody[T any](request *http.Request, decode func(*http.Request, any) error, rules validation.Rules[T]) (T, error) {
	var value T
	if decode == nil { decode = httprequest.DecodeJSON }
	if err := decode(request, &value); err != nil { return value, &DecodeError{Err: err} }
	return validation.Process(value, rules)
}

// Run binds path -> query -> JSON body, normalizes and validates each input,
// then invokes Service once and writes exactly one HTTP response. It deliberately
// does not return a value: net/http handlers return no value.
func Run[Path, Query, Body, Response any](writer http.ResponseWriter, request *http.Request, config Config[Path, Query, Body, Response]) {
	if config.Service == nil {
		writeFailure(writer, request, config.Writer, fmt.Errorf("REST service callback is missing"), config.Fallback, config.Fallback)
		return
	}
	var input Input[Path, Query, Body]
	var err error
	if input.Path, err = bindPath[Path](request); err != nil {
		writeFailure(writer, request, config.Writer, err, config.InvalidPath, config.Fallback)
		return
	}
	if input.Path, err = validation.Process(input.Path, config.Path); err != nil {
		writeFailure(writer, request, config.Writer, err, config.InvalidPath, config.Fallback)
		return
	}
	if input.Query, err = bindQuery[Query](request); err != nil {
		writeFailure(writer, request, config.Writer, err, config.InvalidQuery, config.Fallback)
		return
	}
	if input.Query, err = validation.Process(input.Query, config.Query); err != nil {
		writeFailure(writer, request, config.Writer, err, config.InvalidQuery, config.Fallback)
		return
	}
	if !isEmpty[Body]() {
		if input.Body, err = ReadBody[Body](request, config.Decode, config.Body); err != nil {
			writeFailure(writer, request, config.Writer, err, config.InvalidBody, config.Fallback)
			return
		}
	}
	if config.Validate != nil {
		if err = config.Validate(input); err != nil {
			writeFailure(writer, request, config.Writer, err, config.InvalidBody, config.Fallback)
			return
		}
	}
	result, err := config.Service(request.Context(), input)
	if err != nil {
		writeFailure(writer, request, config.Writer, err, config.Fallback, config.Fallback)
		return
	}
	writeSuccess(writer, request, config.Writer, result, config.Fallback)
}

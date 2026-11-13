package rest

import (
	"errors"
	"fmt"
	"net/http"

	validation "grovehub/libs/core/annotation/validation"
	coreerrors "grovehub/libs/core/errors"
	httpresponse "grovehub/libs/core/http/response"
)

// Failure is the safe, final response descriptor. The raw cause is never
// serialized. Application metadata may be included when explicitly provided.
type Failure struct {
	StatusCode int
	Definition coreerrors.Definition
	Metadata   any
}

// Writer optionally overrides the default GroveHub response envelope.
type Writer[T any] struct {
	Success func(http.ResponseWriter, *http.Request, Result[T]) error
	Error   func(http.ResponseWriter, *http.Request, Failure) error
}

func writeSuccess[T any](writer http.ResponseWriter, request *http.Request, output Writer[T], result Result[T], fallback coreerrors.Definition) {
	status := result.StatusCode
	if status == 0 { status = http.StatusOK }
	if status < 200 || status > 299 {
		writeFailure(writer, request, output, fmt.Errorf("service returned a non-success result status"), fallback, fallback)
		return
	}
	if output.Success != nil {
		_ = output.Success(writer, request, result)
		return
	}
	for name, values := range result.Headers {
		for _, value := range values { writer.Header().Add(name, value) }
	}
	if status == http.StatusNoContent {
		writer.WriteHeader(status)
		return
	}
	_ = httpresponse.WriteSuccess(writer, status, result.Data)
}

func writeFailure[T any](writer http.ResponseWriter, request *http.Request, output Writer[T], err error, invalid, fallback coreerrors.Definition) {
	failure := classifyFailure(err, invalid, fallback)
	if output.Error != nil {
		_ = output.Error(writer, request, failure)
		return
	}
	if failure.Metadata == nil {
		_ = httpresponse.WriteErrorContext(request.Context(), writer, failure.StatusCode, failure.Definition)
		return
	}
	_ = httpresponse.WriteErrorWithMetadataContext(request.Context(), writer, failure.StatusCode, failure.Definition, failure.Metadata)
}

func classifyFailure(err error, invalid, fallback coreerrors.Definition) Failure {
	var appError *coreerrors.Error
	if errors.As(err, &appError) {
		return Failure{StatusCode: appError.Definition().HTTPStatus(), Definition: appError.Definition(), Metadata: appError.Metadata()}
	}
	var decodeError *DecodeError
	var bindError BindError
	var dtoError validation.Error
	if errors.As(err, &decodeError) || errors.As(err, &bindError) || errors.As(err, &dtoError) {
		if invalid.Code() != "" && invalid.HTTPStatus() >= 400 && invalid.HTTPStatus() < 500 {
			return Failure{StatusCode: invalid.HTTPStatus(), Definition: invalid}
		}
	}
	// Never publish raw, unexpected errors or accidentally turn an internal
	// failure into a client error because the fallback was a 4xx definition.
	return Failure{StatusCode: http.StatusInternalServerError, Definition: fallback}
}

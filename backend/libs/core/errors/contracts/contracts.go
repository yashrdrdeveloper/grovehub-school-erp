package contracts

import "net/http"

type Code string

func (code Code) String() string {
	return string(code)
}

type Definition struct {
	code       Code
	message    string
	httpStatus int
}

// NewDefinition declares an application error with an optional HTTP status.
// Older two-argument definitions remain valid and default to HTTP 500.
// Configured statuses must be HTTP 4xx/5xx; invalid declarations fail fast.
func NewDefinition(code Code, message string, httpStatus ...int) Definition {
	status := http.StatusInternalServerError
	if len(httpStatus) > 1 {
		panic("error definition accepts at most one HTTP status")
	}
	if len(httpStatus) == 1 {
		status = httpStatus[0]
		if status < http.StatusBadRequest || status > 599 {
			panic("error definition requires an HTTP 4xx or 5xx status")
		}
	}
	return Definition{code: code, message: message, httpStatus: status}
}

func (definition Definition) Code() Code {
	return definition.code
}

func (definition Definition) Message() string {
	return definition.message
}

// HTTPStatus is the canonical transport status for the error definition.
// Zero-valued definitions and legacy definitions default to HTTP 500.
func (definition Definition) HTTPStatus() int {
	if definition.httpStatus == 0 {
		return http.StatusInternalServerError
	}
	return definition.httpStatus
}

type Error struct {
	definition Definition
	cause      error
	metadata   any
}

func NewError(definition Definition) *Error {
	return &Error{definition: definition}
}

func NewErrorWithMetadata(definition Definition, metadata any) *Error {
	return &Error{definition: definition, metadata: metadata}
}

func Wrap(definition Definition, cause error) *Error {
	return &Error{definition: definition, cause: cause}
}

func WrapWithMetadata(definition Definition, cause error, metadata any) *Error {
	return &Error{definition: definition, cause: cause, metadata: metadata}
}

func (err *Error) Error() string {
	return err.definition.Message()
}

func (err *Error) Unwrap() error {
	return err.cause
}

func (err *Error) Definition() Definition {
	return err.definition
}

func (err *Error) Metadata() any {
	return err.metadata
}

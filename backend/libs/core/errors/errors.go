// Package coreerrors exposes application error definitions and wrapped errors.
package coreerrors

import (
	"grovehub/libs/core/errors/internal/src/contracts"
	"grovehub/libs/core/errors/internal/src/operations/create"
	"grovehub/libs/core/errors/internal/src/operations/definition"
	"grovehub/libs/core/errors/internal/src/operations/wrap"
)

type Code = contracts.Code
type Definition = contracts.Definition
type Error = contracts.Error

func NewDefinition(code Code, message string, httpStatus ...int) Definition {
	return definition.NewDefinition(code, message, httpStatus...)
}

func NewError(definition Definition) *Error {
	return create.NewError(definition)
}

func NewErrorWithMetadata(definition Definition, metadata any) *Error {
	return create.NewErrorWithMetadata(definition, metadata)
}

func Wrap(definition Definition, cause error) *Error {
	return wrap.Wrap(definition, cause)
}

func WrapWithMetadata(definition Definition, cause error, metadata any) *Error {
	return wrap.WrapWithMetadata(definition, cause, metadata)
}

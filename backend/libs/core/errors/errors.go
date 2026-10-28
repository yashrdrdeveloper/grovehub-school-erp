package coreerrors

import (
	"grovehub/libs/core/errors/contracts"
	"grovehub/libs/core/errors/helpers"
)

type Code = contracts.Code
type Definition = contracts.Definition
type Error = contracts.Error

func NewDefinition(code Code, message string) Definition {
	return helpers.NewDefinition(code, message)
}

func NewError(definition Definition) *Error {
	return helpers.NewError(definition)
}

func NewErrorWithMetadata(definition Definition, metadata any) *Error {
	return helpers.NewErrorWithMetadata(definition, metadata)
}

func Wrap(definition Definition, cause error) *Error {
	return helpers.Wrap(definition, cause)
}

func WrapWithMetadata(definition Definition, cause error, metadata any) *Error {
	return helpers.WrapWithMetadata(definition, cause, metadata)
}

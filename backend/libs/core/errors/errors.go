package coreerrors

import (
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors/helpers"
)

type Code = contracts.Code
type Definition = contracts.Definition
type Error = contracts.Error

func NewDefinition(code Code, message string, httpStatus ...int) Definition {
	return helpers.NewDefinition(code, message, httpStatus...)
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

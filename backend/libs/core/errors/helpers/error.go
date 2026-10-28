package helpers

import "grovehub/libs/core/errors/contracts"

func NewDefinition(code contracts.Code, message string) contracts.Definition {
	return contracts.NewDefinition(code, message)
}

func NewError(definition contracts.Definition) *contracts.Error {
	return contracts.NewError(definition)
}

func NewErrorWithMetadata(definition contracts.Definition, metadata any) *contracts.Error {
	return contracts.NewErrorWithMetadata(definition, metadata)
}

func Wrap(definition contracts.Definition, cause error) *contracts.Error {
	return contracts.Wrap(definition, cause)
}

func WrapWithMetadata(definition contracts.Definition, cause error, metadata any) *contracts.Error {
	return contracts.WrapWithMetadata(definition, cause, metadata)
}

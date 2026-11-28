package wrap

import "grovehub/libs/core/errors/internal/src/contracts"

func Wrap(definition contracts.Definition, cause error) *contracts.Error {
	return contracts.NewError(definition, cause, nil)
}

func WrapWithMetadata(definition contracts.Definition, cause error, metadata any) *contracts.Error {
	return contracts.NewError(definition, cause, metadata)
}

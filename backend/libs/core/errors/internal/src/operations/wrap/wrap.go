package wrap

import "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors/internal/src/contracts"

func Wrap(definition contracts.Definition, cause error) *contracts.Error {
	return contracts.NewError(definition, cause, nil)
}

func WrapWithMetadata(definition contracts.Definition, cause error, metadata any) *contracts.Error {
	return contracts.NewError(definition, cause, metadata)
}

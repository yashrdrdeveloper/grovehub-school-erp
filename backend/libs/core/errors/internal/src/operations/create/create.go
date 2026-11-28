package create

import "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors/internal/src/contracts"

func NewError(definition contracts.Definition) *contracts.Error {
	return contracts.NewError(definition, nil, nil)
}

func NewErrorWithMetadata(definition contracts.Definition, metadata any) *contracts.Error {
	return contracts.NewError(definition, nil, metadata)
}

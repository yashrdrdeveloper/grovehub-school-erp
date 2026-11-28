package project

import (
	coreerrors "grovehub/libs/core/errors"
	"grovehub/libs/core/utils"

	definitions "grovehub/services/registry/internal/app/constants/error"
)

func Validate(projectID string) error {
	if utils.ValidateUUID(projectID) != nil {
		return coreerrors.NewError(definitions.Invalid)
	}
	return nil
}

package helpers

import (
	"errors"

	coreerrors "grovehub/libs/core/errors"
	definitions "grovehub/services/registry/internal/app/constants/error"
	"grovehub/services/registry/internal/app/repository/constants"
)

func RepositoryError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, constants.ErrNotFound) {
		return coreerrors.Wrap(definitions.NotFound, err)
	}
	if errors.Is(err, constants.ErrConflict) {
		return coreerrors.Wrap(definitions.Conflict, err)
	}
	return coreerrors.Wrap(definitions.Unavailable, err)
}

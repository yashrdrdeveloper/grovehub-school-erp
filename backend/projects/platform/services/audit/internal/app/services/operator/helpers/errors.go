package helpers

import (
	"errors"

	coreerrors "grovehub/libs/core/errors"
	definitions "grovehub/services/audit/internal/app/constants/error"
	repository "grovehub/services/audit/internal/app/repository/operator"
)

func StorageError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, repository.ErrConflict) {
		return coreerrors.Wrap(definitions.OperatorConflict, err)
	}
	return coreerrors.Wrap(definitions.OperatorUnavailable, err)
}

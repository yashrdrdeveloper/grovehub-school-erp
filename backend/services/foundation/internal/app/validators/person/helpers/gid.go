package helpers

import (
    coreerrors "grovehub/libs/core/errors"
    corenormalize "grovehub/libs/core/normalize"
    personerror "grovehub/services/foundation/internal/app/constants/error"
)

func ValidateGID(gid string) error {
	if _, err := corenormalize.UUIDv7(gid); err != nil {
		return coreerrors.NewError(personerror.InvalidPersonGID)
	}
	return nil
}

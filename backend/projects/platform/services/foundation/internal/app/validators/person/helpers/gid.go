package helpers

import (
    coreerrors "grovehub/libs/core/errors"
    coreutils "grovehub/libs/core/utils"
    personerror "grovehub/services/foundation/internal/app/constants/error"
)

func ValidateGID(gid string) error {
	if _, err := coreutils.NormalizeUUIDv7(gid); err != nil {
		return coreerrors.NewError(personerror.InvalidPersonGID)
	}
	return nil
}

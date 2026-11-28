package helpers

import (
    coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
    coreutils "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/utils"
    personerror "grovehub/services/foundation/internal/app/constants/error"
)

func ValidateGID(gid string) error {
	if _, err := coreutils.NormalizeUUIDv7(gid); err != nil {
		return coreerrors.NewError(personerror.InvalidPersonGID)
	}
	return nil
}

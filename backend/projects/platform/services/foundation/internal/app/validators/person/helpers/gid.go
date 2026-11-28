package helpers

import (
    coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
    corenormalize "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/normalize"
    personerror "grovehub/services/foundation/internal/app/constants/error"
)

func ValidateGID(gid string) error {
	if _, err := corenormalize.UUIDv7(gid); err != nil {
		return coreerrors.NewError(personerror.InvalidPersonGID)
	}
	return nil
}

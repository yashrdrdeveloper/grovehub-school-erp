package helpers

import (
	coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
	"grovehub/services/audit/internal/app/bootstrap/constants"
	"grovehub/services/audit/internal/app/repository/schema"
)

func CheckSchema(repository schema.Repository) error {
	if err := repository.CheckActivity(); err != nil {
		return coreerrors.Wrap(constants.ActivityMigrationRequired, err)
	}

	if err := repository.CheckRejection(); err != nil {
		return coreerrors.Wrap(constants.RejectionMigrationRequired, err)
	}

	if err := repository.CheckReceipt(); err != nil {
		return coreerrors.Wrap(constants.ReceiptMigrationRequired, err)
	}

	if err := repository.CheckOperator(); err != nil {
		return coreerrors.Wrap(constants.OperatorMigrationRequired, err)
	}

	if err := repository.CheckSchoolReceipt(); err != nil {
		return coreerrors.Wrap(constants.SchoolReceiptMigrationRequired, err)
	}

	return nil
}

package impl

import (
	"context"
	"time"

	"gorm.io/gorm"

	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/repository/account"
	"grovehub/services/identity/internal/app/repository/account/constants"
)

type SetupWriter struct {
	db *gorm.DB
}

func NewSetupWriter(db *gorm.DB) account.SetupWriter {
	return &SetupWriter{db: db}
}

func (writer *SetupWriter) ActivateSetup(ctx context.Context, tenantID string, accountID string, version int64, at time.Time) error {
	if writer.db == nil || !transaction.Active(ctx) || tenantID == "" || accountID == "" || version < 1 || at.IsZero() {
		return account.ErrConflict
	}
	result := transaction.Database(ctx, writer.db).Exec(constants.ChangeStatus, entity.Active, at, tenantID, accountID, version, entity.PendingSetup)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return account.ErrVersionConflict
	}
	return nil
}

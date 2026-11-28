package impl

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/repository/account"
	"grovehub/services/identity/internal/app/repository/account/constants"
)

type StatusWriter struct {
	db *gorm.DB
}

var _ account.StatusWriter = (*StatusWriter)(nil)

func NewStatusWriter(db *gorm.DB) account.StatusWriter {
	return &StatusWriter{db: db}
}

func (writer *StatusWriter) ChangeStatus(ctx context.Context, request account.StatusChange) error {
	id, err := uuid.Parse(request.AccountID)
	if err != nil || id == uuid.Nil || strings.TrimSpace(request.TenantID) == "" || strings.TrimSpace(request.TenantID) != request.TenantID {
		return account.ErrNotFound
	}
	if request.Version < 1 || request.UpdatedAt.IsZero() || !validStatusChange(request.Before, request.After) {
		return account.ErrConflict
	}
	if writer.db == nil {
		return transaction.ErrMissingDatabase
	}
	result := transaction.Database(ctx, writer.db).Exec(constants.ChangeStatus, request.After, request.UpdatedAt, request.TenantID, id.String(), request.Version, request.Before)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return account.ErrVersionConflict
	}
	return nil
}

func validStatusChange(before entity.AccountStatus, after entity.AccountStatus) bool {
	switch before {
	case entity.Active:
		return after == entity.Suspended || after == entity.Deactivated
	case entity.Suspended:
		return after == entity.Active || after == entity.Deactivated
	case entity.PendingSetup:
		return after == entity.Deactivated
	default:
		return false
	}
}

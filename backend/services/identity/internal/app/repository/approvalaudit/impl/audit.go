package impl

import (
	"context"

	"gorm.io/gorm"

	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	"grovehub/services/identity/internal/app/repository/approvalaudit"
	"grovehub/services/identity/internal/app/repository/approvalaudit/constants"
)

type Writer struct {
	db *gorm.DB
}

var _ approvalaudit.Writer = (*Writer)(nil)

func New(db *gorm.DB) approvalaudit.Writer {
	return &Writer{db: db}
}

func (writer *Writer) Insert(ctx context.Context, entry approvalaudit.Entry) error {
	if writer.db == nil {
		return transaction.ErrMissingDatabase
	}
	return transaction.Database(ctx, writer.db).Exec(constants.Insert, entry.EventID, entry.TenantID, entry.RequestID, entry.ActorID, entry.Version, entry.Operation, string(entry.Payload), entry.OccurredAt).Error
}

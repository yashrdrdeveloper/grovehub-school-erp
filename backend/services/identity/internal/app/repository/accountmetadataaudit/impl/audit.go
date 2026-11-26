package impl

import (
	"context"

	"gorm.io/gorm"

	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	"grovehub/services/identity/internal/app/repository/accountmetadataaudit"
	"grovehub/services/identity/internal/app/repository/accountmetadataaudit/constants"
)

type Writer struct {
	db *gorm.DB
}

func New(db *gorm.DB) accountmetadataaudit.Writer {
	return &Writer{db: db}
}

func (writer *Writer) Insert(ctx context.Context, entry accountmetadataaudit.Entry) error {
	if writer.db == nil || !transaction.Active(ctx) {
		return transaction.ErrMissingDatabase
	}
	return transaction.Database(ctx, writer.db).Exec(constants.Insert, entry.EventID, entry.TenantID, entry.ActorID, entry.AccountID, entry.Version, string(entry.Payload), entry.OccurredAt).Error
}

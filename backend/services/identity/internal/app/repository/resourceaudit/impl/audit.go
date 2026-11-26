package impl

import (
	"context"
	"errors"

	"gorm.io/gorm"

	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	"grovehub/services/identity/internal/app/repository/resourceaudit"
	"grovehub/services/identity/internal/app/repository/resourceaudit/constants"
)

type Writer struct {
	db *gorm.DB
}

var _ resourceaudit.Writer = (*Writer)(nil)

func New(db *gorm.DB) resourceaudit.Writer {
	return &Writer{db: db}
}

func (writer *Writer) Insert(ctx context.Context, entry resourceaudit.Entry) error {
	if writer.db == nil || !transaction.Active(ctx) {
		return errors.New("resource audit requires active transaction")
	}
	return transaction.Database(ctx, writer.db).Exec(constants.Insert, entry.EventID, entry.TenantID, entry.RequestID, entry.ActorID, entry.ResourceType, entry.ResourceID, entry.Version, entry.Operation, string(entry.Payload), entry.OccurredAt).Error
}

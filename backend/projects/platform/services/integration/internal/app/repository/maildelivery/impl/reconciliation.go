package impl

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	coredb "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
	"grovehub/services/integration/internal/app/repository/maildelivery"
)

type receiptReader struct { db *gorm.DB }

func NewReceiptReader(db *gorm.DB) maildelivery.ReceiptReader {
	return &receiptReader{db:db}
}

func (r *receiptReader) ListReceiptPage(ctx context.Context, tenantID, after string, limit int) ([]maildelivery.Receipt,error) {
	id, err := uuid.Parse(after)
	if r.db == nil || tenantID == "" || err != nil || id.String() != after || limit < 1 || limit > 500 {
		return nil, errors.New("invalid receipt page")
	}
	rows := make([]maildelivery.Receipt,0)
	err = coredb.TransactionDB(ctx,r.db).Raw(`
SELECT delivery_id::text AS delivery_id, state, reserved_at, accepted_at
FROM integration.mail_deliveries
WHERE tenant_id = ? AND delivery_id > CAST(? AS uuid)
ORDER BY delivery_id LIMIT ?`,tenantID,after,limit).Scan(&rows).Error
	return rows,err
}

package impl

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	coredb "grovehub/libs/core/database"
	"grovehub/services/identity/client/invitation"
)

type deliveryReader struct { db *gorm.DB }

func NewDeliveryReader(db *gorm.DB) invitation.DeliveryReader {
	return &deliveryReader{db: db}
}

func (r *deliveryReader) ListDeliveryPage(ctx context.Context, tenantID, after string, limit int) ([]invitation.Delivery, error) {
	id, err := uuid.Parse(after)
	if r.db == nil || tenantID == "" || err != nil || id.String() != after || limit < 1 || limit > 500 {
		return nil, errors.New("invalid invitation delivery page")
	}
	rows := make([]invitation.Delivery, 0)
	err = coredb.TransactionDB(ctx, r.db).Raw(`
SELECT i.id::text AS delivery_id, i.account_id::text AS account_id,
       CASE WHEN i.purpose IN ('ACCOUNT_SETUP','PASSWORD_RESET') THEN i.purpose ELSE 'UNKNOWN' END AS purpose,
       i.sent_at, i.suppressed_at, i.suppression_reason, i.lease_until
FROM identity.invitations i JOIN identity.accounts a ON a.account_id = i.account_id
WHERE a.tenant_id = ? AND i.id > CAST(? AS uuid)
ORDER BY i.id LIMIT ?`, tenantID, after, limit).Scan(&rows).Error
	return rows, err
}

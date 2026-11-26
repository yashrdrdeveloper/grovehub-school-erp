package impl

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"grovehub/services/integration/internal/app/repository/maildelivery"
)

// These writes intentionally do not join a caller's transaction. A provider
// dispatch reservation must survive rollback of Identity's delivery update.
type Repository struct { db *gorm.DB }

func New(db *gorm.DB) maildelivery.Repository {
	return &Repository{db: db}
}

func (repository *Repository) Find(ctx context.Context, tenantID string, deliveryID string) (maildelivery.Record, bool, error) {
	var record maildelivery.Record
	result := repository.db.WithContext(ctx).Raw("SELECT request_hash, state FROM integration.mail_deliveries WHERE tenant_id = ? AND delivery_id = ?", tenantID, deliveryID).Scan(&record)
	return record, result.RowsAffected == 1, result.Error
}

func (repository *Repository) Reserve(ctx context.Context, tenantID string, deliveryID string, requestHash string) (bool, error) {
	result := repository.db.WithContext(ctx).Exec("INSERT INTO integration.mail_deliveries(tenant_id, delivery_id, request_hash, state) VALUES (?, ?, ?, 'DISPATCHING') ON CONFLICT (tenant_id, delivery_id) DO NOTHING", tenantID, deliveryID, requestHash)
	return result.RowsAffected == 1, result.Error
}

func (repository *Repository) Accept(ctx context.Context, tenantID string, deliveryID string, requestHash string) error {
	result := repository.db.WithContext(ctx).Exec("UPDATE integration.mail_deliveries SET state = 'ACCEPTED', accepted_at = clock_timestamp() WHERE tenant_id = ? AND delivery_id = ? AND request_hash = ? AND state = 'DISPATCHING'", tenantID, deliveryID, requestHash)
	if result.Error != nil { return result.Error }
	if result.RowsAffected != 1 { return errors.New("mail dispatch reservation missing") }
	return nil
}

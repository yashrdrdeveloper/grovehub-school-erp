package impl

import (
	"context"

	"gorm.io/gorm"

	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	repo "grovehub/services/identity/internal/app/repository/account"
	"grovehub/services/identity/internal/app/repository/account/constants"
)

type RegistrationWriter struct {
	db *gorm.DB
}

var _ repo.RegistrationWriter = (*RegistrationWriter)(nil)

func NewRegistrationWriter(db *gorm.DB) repo.RegistrationWriter {
	return &RegistrationWriter{db: db}
}

func (writer *RegistrationWriter) FindRegistration(ctx context.Context, tenantID string, actorID string, key string) (repo.RegistrationReceipt, error) {
	if writer.db == nil || !transaction.Active(ctx) || tenantID == "" || actorID == "" || key == "" {
		return repo.RegistrationReceipt{}, repo.ErrConflict
	}
	var receipt repo.RegistrationReceipt
	result := transaction.Database(ctx, writer.db).Raw(constants.FindRegistration, tenantID, actorID, key).Scan(&receipt)
	if result.Error != nil {
		return repo.RegistrationReceipt{}, result.Error
	}
	if result.RowsAffected != 1 {
		return repo.RegistrationReceipt{}, repo.ErrNotFound
	}
	return receipt, nil
}

func (writer *RegistrationWriter) CreatePending(ctx context.Context, registration repo.Registration) error {
	if writer.db == nil || !transaction.Active(ctx) || registration.TenantID == "" {
		return repo.ErrConflict
	}
	db := transaction.Database(ctx, writer.db)
	if err := db.Exec(constants.InsertPending, registration.AccountID, registration.PersonGID, registration.TenantID, registration.Email, registration.Phone, registration.Email, registration.At, registration.At).Error; err != nil {
		return mapWriteError(err)
	}
	return mapWriteError(db.Exec(constants.InsertRegistrationReceipt, registration.TenantID, registration.ActorID, registration.Key, registration.Digest, registration.AccountID, registration.EventID, registration.At).Error)
}

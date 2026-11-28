package impl

import (
	"context"

	"gorm.io/gorm"

	coredatabase "grovehub/libs/core/database"
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
	if writer.db == nil || !coredatabase.TransactionActive(ctx) || tenantID == "" || actorID == "" || key == "" {
		return repo.RegistrationReceipt{}, repo.ErrConflict
	}
	var receipt repo.RegistrationReceipt
	result := coredatabase.TransactionDB(ctx, writer.db).Raw(constants.FindRegistration, tenantID, actorID, key).Scan(&receipt)
	if result.Error != nil {
		return repo.RegistrationReceipt{}, result.Error
	}
	if result.RowsAffected != 1 {
		return repo.RegistrationReceipt{}, repo.ErrNotFound
	}
	return receipt, nil
}

func (writer *RegistrationWriter) CreatePending(ctx context.Context, registration repo.Registration) error {
	if writer.db == nil || !coredatabase.TransactionActive(ctx) || registration.TenantID == "" {
		return repo.ErrConflict
	}
	db := coredatabase.TransactionDB(ctx, writer.db)
	if err := db.Exec(constants.InsertPending, registration.AccountID, registration.PersonGID, registration.TenantID, registration.Email, registration.Phone, registration.Email, registration.At, registration.At).Error; err != nil {
		return mapWriteError(err)
	}
	return mapWriteError(db.Exec(constants.InsertRegistrationReceipt, registration.TenantID, registration.ActorID, registration.Key, registration.Digest, registration.AccountID, registration.EventID, registration.At).Error)
}

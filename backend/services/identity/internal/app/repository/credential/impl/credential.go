package impl

import (
	"context"
	"errors"

	"gorm.io/gorm"

	transaction "grovehub/libs/core/database/transaction"
	"grovehub/services/identity/internal/app/entity"
	credentialrepo "grovehub/services/identity/internal/app/repository/credential"
)

type Repository struct {
	db *gorm.DB
}

var _ credentialrepo.Repository = (*Repository)(nil)

func New(db *gorm.DB) credentialrepo.Repository {
	return &Repository{db: db}
}

func (repository *Repository) Find(ctx context.Context, accountID string) (entity.PasswordCredential, error) {
	var credential entity.PasswordCredential
	err := transaction.Database(ctx, repository.db).
		Where("account_id = ?", accountID).
		Take(&credential).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.PasswordCredential{}, credentialrepo.ErrNotFound
	}
	return credential, err
}

func (repository *Repository) Create(ctx context.Context, credential entity.PasswordCredential) error {
	return transaction.Database(ctx, repository.db).Create(&credential).Error
}

func (repository *Repository) Update(ctx context.Context, credential entity.PasswordCredential) error {
	result := transaction.Database(ctx, repository.db).
		Model(&entity.PasswordCredential{}).
		Where("account_id = ? AND version = ?", credential.AccountID, credential.Version-1).
		Updates(map[string]any{
			"password_hash": credential.PasswordHash,
			"password_changed_at": credential.PasswordChangedAt,
			"version": credential.Version,
			"updated_at": credential.UpdatedAt,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return errors.New("credential version conflict")
	}
	return nil
}

package impl

import (
	"context"
	"time"

	"gorm.io/gorm"

	transaction "grovehub/libs/core/database/transaction"
	"grovehub/services/identity/internal/app/entity"
	tokenrepo "grovehub/services/identity/internal/app/repository/passwordtoken"
)

type Repository struct {
	db *gorm.DB
}

var _ tokenrepo.Repository = (*Repository)(nil)

func New(db *gorm.DB) tokenrepo.Repository {
	return &Repository{db: db}
}

func (repository *Repository) Issue(ctx context.Context, token entity.PasswordToken) error {
	return transaction.Database(ctx, repository.db).Create(&token).Error
}

func (repository *Repository) Consume(ctx context.Context, hash string, purpose entity.TokenPurpose, now time.Time) (string, error) {
	var result struct {
		AccountID string `gorm:"column:account_id"`
	}
	query := transaction.Database(ctx, repository.db).Raw(`
		UPDATE identity.password_tokens
		SET used_at = ?
		WHERE token_hash = ?
		  AND purpose = ?
		  AND used_at IS NULL
		  AND expires_at > ?
		RETURNING account_id
	`, now, hash, purpose, now).Scan(&result)
	if query.Error != nil {
		return "", query.Error
	}
	if result.AccountID == "" {
		return "", tokenrepo.ErrInvalidToken
	}
	return result.AccountID, nil
}

func (repository *Repository) Invalidate(ctx context.Context, accountID string, purpose entity.TokenPurpose) error {
	return transaction.Database(ctx, repository.db).
		Model(&entity.PasswordToken{}).
		Where("account_id = ? AND purpose = ? AND used_at IS NULL", accountID, purpose).
		Update("used_at", time.Now().UTC()).Error
}

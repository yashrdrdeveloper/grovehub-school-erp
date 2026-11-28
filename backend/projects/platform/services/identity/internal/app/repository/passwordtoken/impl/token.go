package impl

import (
	"context"
	"time"

	"gorm.io/gorm"

	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
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
	return coredatabase.TransactionDB(ctx, repository.db).Create(&token).Error
}

func (repository *Repository) Consume(ctx context.Context, hash string, purpose entity.TokenPurpose, now time.Time) (string, error) {
	var result struct {
		AccountID string `gorm:"column:account_id"`
	}
	query := coredatabase.TransactionDB(ctx, repository.db).Raw(`
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


func (repository *Repository) Inspect(ctx context.Context, hash string, purpose entity.TokenPurpose, now time.Time) (entity.PasswordToken, error) {
  var token entity.PasswordToken
  err := coredatabase.TransactionDB(ctx, repository.db).
    Where("token_hash = ? AND purpose = ? AND used_at IS NULL AND expires_at > ?", hash, purpose, now).
    Take(&token).Error
  if err == gorm.ErrRecordNotFound { return entity.PasswordToken{}, tokenrepo.ErrInvalidToken }
  return token, err
}

// ReserveSetupAttempt commits outside the setup unit of work. A failed password
// update must not roll back this counter, or retries would be unlimited.
func (repository *Repository) ReserveSetupAttempt(ctx context.Context, hash string, now time.Time) error {
  result := repository.db.WithContext(ctx).Exec(`
    UPDATE identity.password_tokens SET attempts = attempts + 1
    WHERE token_hash = ? AND purpose = ? AND used_at IS NULL AND expires_at > ? AND attempts < 2
  `, hash, entity.SetupGrant, now)
  if result.Error != nil { return result.Error }
  if result.RowsAffected == 0 { return tokenrepo.ErrAttemptsExceeded }
  return nil
}

func (repository *Repository) Invalidate(ctx context.Context, accountID string, purpose entity.TokenPurpose) error {
	return coredatabase.TransactionDB(ctx, repository.db).
		Model(&entity.PasswordToken{}).
		Where("account_id = ? AND purpose = ? AND used_at IS NULL", accountID, purpose).
		Update("used_at", time.Now().UTC()).Error
}

// ReserveResetAttempt persists the attempt before the credential transaction.
// Failed password replacement must not restore the attempt budget.
func (repository *Repository) ReserveResetAttempt(ctx context.Context, hash string, now time.Time) error {
	result := repository.db.WithContext(ctx).Exec(`
		UPDATE identity.password_tokens SET attempts = attempts + 1
		WHERE token_hash = ? AND purpose = ? AND used_at IS NULL
		  AND expires_at > ? AND attempts < 2 AND account_version IS NOT NULL
	`, hash, entity.ResetGrant, now)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 1 {
		return nil
	}
	proof, err := repository.Inspect(ctx, hash, entity.ResetGrant, time.Now().UTC())
	if err != nil {
		return err
	}
	if proof.AccountVersion == nil {
		return tokenrepo.ErrInvalidToken
	}
	if proof.Attempts >= 2 {
		return tokenrepo.ErrResetAttemptsExceeded
	}
	return tokenrepo.ErrInvalidToken
}

package impl

import (
    "context"
    "errors"
    "time"

    "gorm.io/gorm"
    coretransaction "grovehub/libs/core/database/transaction"
    "grovehub/services/identity/entity"
    identityrepository "grovehub/services/identity/repository"
)

type PasswordSetupTokenRepository struct { db *gorm.DB }

func NewPasswordSetupTokenRepository(db *gorm.DB) identityrepository.PasswordSetupTokenRepository {
    return &PasswordSetupTokenRepository{db: db}
}

func (repository *PasswordSetupTokenRepository) Issue(ctx context.Context, token entity.PasswordSetupToken) error {
    return coretransaction.Database(ctx, repository.db).Create(&token).Error
}

func (repository *PasswordSetupTokenRepository) Consume(ctx context.Context, tokenHash string, now time.Time) (entity.AccountId, error) {
    var row struct { AccountID string `gorm:"column:account_id"` }
    err := coretransaction.Database(ctx, repository.db).Raw(
        "UPDATE identity.password_setup_tokens SET used_at = ? WHERE token_hash = ? AND used_at IS NULL AND expires_at > ? RETURNING account_id",
        now, tokenHash, now,
    ).Scan(&row).Error
    if err != nil { return "", err }
    if row.AccountID == "" { return "", identityrepository.ErrPasswordSetupTokenInvalid }
    return entity.AccountId(row.AccountID), nil
}

// Keep this marker to avoid accidental use of a non-transactional scan for consumption.
var _ = errors.Is

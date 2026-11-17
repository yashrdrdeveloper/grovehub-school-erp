package impl

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	transaction "grovehub/libs/core/database/transaction"
	"grovehub/services/identity/internal/app/entity"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
)

type Repository struct {
	db *gorm.DB
}

var _ sessionrepo.Repository = (*Repository)(nil)

func New(db *gorm.DB) sessionrepo.Repository {
	return &Repository{db: db}
}

func (repository *Repository) Create(ctx context.Context, session entity.Session) error {
	return transaction.Database(ctx, repository.db).Create(&session).Error
}

func (repository *Repository) FindAccess(ctx context.Context, digest string, now time.Time) (entity.Session, error) {
	var session entity.Session
	err := transaction.Database(ctx, repository.db).
		Where("access_hash = ? AND revoked_at IS NULL AND access_expires_at > ?", digest, now).
		Take(&session).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Session{}, sessionrepo.ErrInvalidSession
	}
	return session, err
}

func (repository *Repository) FindRefresh(ctx context.Context, digest string, now time.Time) (entity.Session, error) {
	var session entity.Session
	err := transaction.Database(ctx, repository.db).
		Where("refresh_hash = ? AND revoked_at IS NULL AND refresh_expires_at > ?", digest, now).
		Take(&session).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Session{}, sessionrepo.ErrInvalidSession
	}
	return session, err
}

func (repository *Repository) Rotate(ctx context.Context, session entity.Session, previousRefreshDigest string) error {
	result := transaction.Database(ctx, repository.db).
		Model(&entity.Session{}).
		Where("session_id = ? AND refresh_hash = ? AND revoked_at IS NULL AND refresh_expires_at > ?", session.SessionID, previousRefreshDigest, time.Now().UTC()).
		Updates(map[string]any{
			"access_hash": session.AccessHash,
			"refresh_hash": session.RefreshHash,
			"access_expires_at": session.AccessExpiresAt,
			"updated_at": session.UpdatedAt,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return sessionrepo.ErrInvalidSession
	}
	return nil
}

func (repository *Repository) Reauthenticate(ctx context.Context, sessionID string, when time.Time) error {
	result := transaction.Database(ctx, repository.db).
		Model(&entity.Session{}).
		Where("session_id = ? AND revoked_at IS NULL", sessionID).
		Update("reauthenticated_at", when)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return sessionrepo.ErrInvalidSession
	}
	return nil
}

func (repository *Repository) List(ctx context.Context, accountID string) ([]entity.Session, error) {
	var sessions []entity.Session
	err := transaction.Database(ctx, repository.db).
		Where("account_id = ? AND revoked_at IS NULL AND refresh_expires_at > ?", accountID, time.Now().UTC()).
		Order("created_at DESC").
		Limit(100).
		Find(&sessions).Error
	return sessions, err
}

func (repository *Repository) Revoke(ctx context.Context, accountID string, sessionID string) error {
	result := transaction.Database(ctx, repository.db).
		Model(&entity.Session{}).
		Where("account_id = ? AND session_id = ? AND revoked_at IS NULL", accountID, sessionID).
		Update("revoked_at", time.Now().UTC())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return sessionrepo.ErrInvalidSession
	}
	return nil
}

func (repository *Repository) RevokeAll(ctx context.Context, accountID string) (int64, error) {
	result := transaction.Database(ctx, repository.db).
		Model(&entity.Session{}).
		Where("account_id = ? AND revoked_at IS NULL", accountID).
		Update("revoked_at", time.Now().UTC())
	return result.RowsAffected, result.Error
}

package impl

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
	schoolgateway "grovehub/libs/platform/session"
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
	return coredatabase.RunTransaction(ctx, repository.db, func(tx context.Context) error {
		db := coredatabase.TransactionDB(tx, repository.db)
		// Serialize session creation with account status changes. A login that
		// wins first is included in revocation; a later login sees inactive status.
		var account struct {
			Status entity.AccountStatus
		}
		result := db.Raw("SELECT status FROM identity.accounts WHERE account_id = ? FOR SHARE", session.AccountID).Scan(&account)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 || account.Status != entity.Active {
			return sessionrepo.ErrInvalidSession
		}
		return db.Create(&session).Error
	})
}

func (repository *Repository) FindAccess(ctx context.Context, digest string, now time.Time) (entity.Session, error) {
	var session entity.Session
	err := coredatabase.TransactionDB(ctx, repository.db).
		Where("access_hash = ? AND revoked_at IS NULL AND access_expires_at > ?", digest, now).
		Take(&session).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Session{}, sessionrepo.ErrInvalidSession
	}
	return session, err
}

func (repository *Repository) FindRefresh(ctx context.Context, digest string, now time.Time) (entity.Session, error) {
	var session entity.Session
	err := coredatabase.TransactionDB(ctx, repository.db).
		Where("refresh_hash = ? AND revoked_at IS NULL AND refresh_expires_at > ?", digest, now).
		Take(&session).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Session{}, sessionrepo.ErrInvalidSession
	}
	return session, err
}

func (repository *Repository) Rotate(ctx context.Context, session entity.Session, previousRefreshDigest string) error {
	result := coredatabase.TransactionDB(ctx, repository.db).
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
	result := coredatabase.TransactionDB(ctx, repository.db).
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
	query := coredatabase.TransactionDB(ctx, repository.db).
		Where("account_id = ? AND revoked_at IS NULL AND refresh_expires_at > ?", accountID, time.Now().UTC())
	if actor, ok := schoolgateway.PrincipalFrom(ctx); ok && actor.AccountID == accountID && actor.SessionGeneration > 0 {
		query = query.Where("session_generation = ?", actor.SessionGeneration)
	}
	err := query.Order("created_at DESC").
		Limit(100).
		Find(&sessions).Error
	return sessions, err
}

func (repository *Repository) Revoke(ctx context.Context, accountID string, sessionID string) error {
	result := coredatabase.TransactionDB(ctx, repository.db).
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
	result := coredatabase.TransactionDB(ctx, repository.db).
		Model(&entity.Session{}).
		Where("account_id = ? AND revoked_at IS NULL", accountID).
		Update("revoked_at", time.Now().UTC())
	return result.RowsAffected, result.Error
}

// LockActive validates the exact session and keeps its row locked until the
// caller's credential transaction completes. It must not run outside a transaction.
func (repository *Repository) LockActive(ctx context.Context, accountID string, sessionID string, now time.Time) (entity.Session, error) {
	if !coredatabase.TransactionActive(ctx) {
		return entity.Session{}, coredatabase.ErrMissingDatabase
	}
	var session entity.Session
	err := coredatabase.TransactionDB(ctx, repository.db).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("account_id = ? AND session_id = ? AND revoked_at IS NULL AND access_expires_at > ? AND refresh_expires_at > ?", accountID, sessionID, now, now).
		Take(&session).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Session{}, sessionrepo.ErrInvalidSession
	}
	return session, err
}

 // FindDelegated validates the referenced session without exposing session secrets.
func (repository *Repository) FindDelegated(ctx context.Context, sessionID string, now time.Time) (entity.Session, error) {
 var session entity.Session
 err := coredatabase.TransactionDB(ctx, repository.db).
  Where("session_id = ? AND revoked_at IS NULL AND access_expires_at > ? AND refresh_expires_at > ?", sessionID, now, now).
  Take(&session).Error
 if errors.Is(err, gorm.ErrRecordNotFound) { return entity.Session{}, sessionrepo.ErrInvalidSession }
 return session, err
}

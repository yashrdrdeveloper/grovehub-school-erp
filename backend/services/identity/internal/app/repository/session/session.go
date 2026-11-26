package session

import (
	"context"
	"errors"
	"time"

	"grovehub/services/identity/internal/app/entity"
)

var ErrInvalidSession = errors.New("invalid or expired session")

type Repository interface {
	Create(ctx context.Context, session entity.Session) error
	FindAccess(ctx context.Context, digest string, now time.Time) (entity.Session, error)
	LockActive(ctx context.Context, accountID string, sessionID string, now time.Time) (entity.Session, error)
	FindRefresh(ctx context.Context, digest string, now time.Time) (entity.Session, error)
	Rotate(ctx context.Context, session entity.Session, previousRefreshDigest string) error
	Reauthenticate(ctx context.Context, sessionID string, when time.Time) error
	List(ctx context.Context, accountID string) ([]entity.Session, error)
	Revoke(ctx context.Context, accountID string, sessionID string) error
	RevokeAll(ctx context.Context, accountID string) (int64, error)
}

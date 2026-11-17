package passwordtoken

import (
	"context"
	"errors"
	"time"

	"grovehub/services/identity/internal/app/entity"
)

var ErrInvalidToken = errors.New("invalid or expired token")

type Repository interface {
	Issue(ctx context.Context, token entity.PasswordToken) error
	Consume(ctx context.Context, hash string, purpose entity.TokenPurpose, now time.Time) (string, error)
	Invalidate(ctx context.Context, accountID string, purpose entity.TokenPurpose) error
}

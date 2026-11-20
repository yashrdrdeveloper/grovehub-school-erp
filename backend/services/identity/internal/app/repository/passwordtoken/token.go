package passwordtoken

import (
	"context"
	"errors"
	"time"

	"grovehub/services/identity/internal/app/entity"
)

var ErrInvalidToken = errors.New("invalid or expired token")
var ErrAttemptsExceeded = errors.New("setup attempts exhausted")

type Repository interface {
	Issue(ctx context.Context, token entity.PasswordToken) error
	Consume(ctx context.Context, hash string, purpose entity.TokenPurpose, now time.Time) (string, error)
	Inspect(ctx context.Context, hash string, purpose entity.TokenPurpose, now time.Time) (entity.PasswordToken, error)
	ReserveSetupAttempt(ctx context.Context, hash string, now time.Time) error
	Invalidate(ctx context.Context, accountID string, purpose entity.TokenPurpose) error
}

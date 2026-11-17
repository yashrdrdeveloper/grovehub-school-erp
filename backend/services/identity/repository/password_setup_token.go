package repository

import (
    "context"
    "errors"
    "time"
    "grovehub/services/identity/entity"
)

var ErrPasswordSetupTokenInvalid = errors.New("password setup token is invalid or expired")

// Consume must be one atomic conditional UPDATE ... RETURNING; callers run it
// in the same transaction as credential creation and account activation.
type PasswordSetupTokenRepository interface {
    Issue(ctx context.Context, token entity.PasswordSetupToken) error
    Consume(ctx context.Context, tokenHash string, now time.Time) (entity.AccountId, error)
}

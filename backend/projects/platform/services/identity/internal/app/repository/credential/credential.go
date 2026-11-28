package credential

import (
	"context"
	"errors"

	"grovehub/services/identity/internal/app/entity"
)

var ErrNotFound = errors.New("credential not found")

// Repository persists password credentials. It never exposes hashes to DTOs.
type Repository interface {
	Find(ctx context.Context, accountID string) (entity.PasswordCredential, error)
	Create(ctx context.Context, credential entity.PasswordCredential) error
	Update(ctx context.Context, credential entity.PasswordCredential) error
}

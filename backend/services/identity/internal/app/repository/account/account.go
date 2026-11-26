package account

import (
	"context"
	"errors"

	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
	"grovehub/services/identity/internal/app/entity"
)

var (
	ErrNotFound = errors.New("account not found")
	ErrConflict = errors.New("account conflict")
	ErrVersionConflict = errors.New("account version conflict")
)

// @Repository entity=Account
type Repository interface {
	Create(ctx context.Context, account entity.Account) error
	FindByID(ctx context.Context, accountID string) (entity.Account, error)
	FindByEmail(ctx context.Context, email string) (entity.Account, error)
	FindByUsername(ctx context.Context, username string) (entity.Account, error)
	FindByPersonGID(ctx context.Context, personGID string) (entity.Account, error)
	FindByIdempotencyKey(ctx context.Context, key string) (entity.Account, error)
	Update(ctx context.Context, account entity.Account) error
	Search(ctx context.Context, tenantID string, query coresearch.Query[string]) (coresearch.QueryResult[entity.Account], error)
}

// Reader is the account projection boundary. Normalized readers never write legacy columns.
type Reader interface {
	FindByID(ctx context.Context, accountID string) (entity.Account, error)
	Search(ctx context.Context, tenantID string, query coresearch.Query[string]) (coresearch.QueryResult[entity.Account], error)
}

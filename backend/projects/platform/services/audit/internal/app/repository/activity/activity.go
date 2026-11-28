package activity

import (
	"context"
	"errors"

	coresearch "grovehub/libs/core/search"
	"grovehub/services/audit/internal/app/entity"
)

var ErrConflict = errors.New("audit identity conflict")

// @Repository entity=Activity
type Repository interface {
	Insert(ctx context.Context, activity entity.Activity) error
	Search(ctx context.Context, tenantID string, selfAccountID string, query coresearch.Query[string]) (coresearch.QueryResult[entity.Activity], error)
}

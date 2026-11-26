package activity

import (
	"context"
	"errors"

	dto "grovehub/services/audit/internal/app/dto/activity"
	"grovehub/services/audit/internal/app/entity"
)

var ErrConflict = errors.New("audit identity conflict")

// @Repository entity=Activity
type Repository interface {
	Insert(ctx context.Context, activity entity.Activity) error
	Search(ctx context.Context, tenantID string, selfAccountID string, query dto.Search) ([]entity.Activity, error)
}

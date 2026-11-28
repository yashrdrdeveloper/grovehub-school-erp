package operator

import (
	"context"
	"errors"

	"grovehub/services/audit/internal/app/entity"
)

var ErrConflict = errors.New("operator occurrence conflict")

// Repository owns immutable writes, existence checks and bounded reads.
type Repository interface {
	InsertActivity(context.Context, entity.OperatorActivity) error
	InsertReceipt(context.Context, entity.SchoolMutationReceipt) error
	SchoolExists(context.Context, string) (bool, error)
	SearchActivities(context.Context, int, int) ([]entity.OperatorActivity, error)
	SearchReceipts(context.Context, string, string, int, int) ([]entity.SchoolMutationReceipt, error)
}

package operator

import (
	"context"

	"grovehub/services/audit/internal/app/entity"
	"grovehub/services/audit/internal/app/repository/operator/constants"
)

var ErrConflict = constants.ErrConflict

// Repository owns immutable writes, existence checks and bounded reads.
type Repository interface {
	InsertActivity(context.Context, entity.OperatorActivity) error
	InsertReceipt(context.Context, entity.SchoolMutationReceipt) error
	SchoolExists(context.Context, string) (bool, error)
	SearchActivities(context.Context, int, int) ([]entity.OperatorActivity, error)
	SearchReceipts(context.Context, string, string, int, int) ([]entity.SchoolMutationReceipt, error)
}

package detachedcontext

import (
	"context"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/transaction/helpers"
)

func DetachedContext(ctx context.Context) context.Context {
	return helpers.DetachedContext(ctx)
}

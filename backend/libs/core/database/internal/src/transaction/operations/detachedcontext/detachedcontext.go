package detachedcontext

import (
	"context"

	"grovehub/libs/core/database/internal/src/transaction/helpers"
)

func DetachedContext(ctx context.Context) context.Context {
	return helpers.DetachedContext(ctx)
}

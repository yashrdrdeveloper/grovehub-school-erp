package active

import (
	"context"

	"grovehub/libs/core/database/internal/src/transaction/helpers"
)

func Active(ctx context.Context) bool {
	_, ok := helpers.FromContext(ctx)

	return ok
}

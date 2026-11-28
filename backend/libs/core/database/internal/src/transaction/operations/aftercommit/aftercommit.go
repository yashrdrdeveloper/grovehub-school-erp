package aftercommit

import (
	"context"

	"grovehub/libs/core/database/internal/src/transaction/helpers"
)

func AfterCommit(ctx context.Context, callback func()) bool {
	state, ok := helpers.FromContext(ctx)
	if !ok {
		return false
	}

	state.AddCallback(callback)

	return true
}

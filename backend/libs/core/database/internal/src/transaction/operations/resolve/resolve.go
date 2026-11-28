package resolve

import (
	"context"

	"gorm.io/gorm"

	"grovehub/libs/core/database/internal/src/transaction/helpers"
)

// Database returns the context-bound transaction or the supplied fallback.
func Database(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	if state, ok := helpers.FromContext(ctx); ok {
		return state.Database(ctx)
	}

	return fallback.WithContext(ctx)
}

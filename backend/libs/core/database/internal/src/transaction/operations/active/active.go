package active

import (
	"context"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/transaction/helpers"
)

func Active(ctx context.Context) bool {
	_, ok := helpers.FromContext(ctx)

	return ok
}

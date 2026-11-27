package approval

import (
	"context"

	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
	dto "grovehub/services/identity/internal/app/dto/approval"
)

type Reader interface {
	History(ctx context.Context, request coresearch.Request) (coresearch.Result[dto.Summary], error)
	Inbox(ctx context.Context, request coresearch.Request) (coresearch.Result[dto.Summary], error)
	PendingCount(ctx context.Context) (int64, error)
	Detail(ctx context.Context, requestID string) (dto.Detail, error)
}

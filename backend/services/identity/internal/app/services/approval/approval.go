package approval

import (
	"context"

	dto "grovehub/services/identity/internal/app/dto/approval"
	"grovehub/services/identity/internal/app/entity"
)

// Service is staged and must not be wired to HTTP until effect execution and
// visibility policy are complete. Only registered typed intents are accepted.
type Service interface {
	CreateLifecycle(ctx context.Context, request dto.CreateLifecycle) (entity.ApprovalRequest, error)
	CreatePermissions(ctx context.Context, request dto.CreatePermissions) (entity.ApprovalRequest, error)
	CreateAssignment(ctx context.Context, request dto.CreateAssignment) (entity.ApprovalRequest, error)
	Decide(ctx context.Context, request dto.Review) (entity.ApprovalRequest, error)
	Cancel(ctx context.Context, requestID string, version int64) (entity.ApprovalRequest, error)
}

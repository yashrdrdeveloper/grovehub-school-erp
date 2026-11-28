package approval

import (
	"context"
	"time"

	coresearch "grovehub/libs/core/search"
	"grovehub/services/identity/internal/app/constants/permission"
	dto "grovehub/services/identity/internal/app/dto/approval"
	"grovehub/services/identity/internal/app/entity"
)

type EligibleGrant struct {
	RoleID string `json:"role_id"`
	Permission permission.Code `json:"permission_code"`
}

// Visibility is constructed exclusively from a verified principal and live
// database-loaded grants. No HTTP input may populate this structure.
type Visibility struct {
	TenantID string
	AccountID string
	Now time.Time
	Grants []EligibleGrant
	Inbox bool
}

type Reader interface {
	Search(ctx context.Context, visibility Visibility, query coresearch.Query[string]) (coresearch.QueryResult[dto.Summary], error)
	CountInbox(ctx context.Context, visibility Visibility) (int64, error)
	Get(ctx context.Context, visibility Visibility, requestID string) (entity.ApprovalRequest, error)
}

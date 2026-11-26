package role

import (
	"context"

	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"

	roledto "grovehub/services/identity/internal/app/dto/role"
)

// Service is internal until approval, audit and route enforcement are ready.
type Service interface {
	CreateDraft(ctx context.Context, request roledto.CreateDraft) (roledto.Detail, error)
	Get(ctx context.Context, roleID string) (roledto.Detail, error)
	UpdateMetadata(ctx context.Context, roleID string, request roledto.UpdateMetadata) (roledto.Detail, error)
	Search(ctx context.Context, request coresearch.Request) (coresearch.Result[roledto.Detail], error)
	Permissions(ctx context.Context, roleID string) (roledto.PermissionSet, error)
	PreviewPermissions(ctx context.Context, roleID string, request roledto.PreviewPermissions) (roledto.PermissionPreview, error)
}

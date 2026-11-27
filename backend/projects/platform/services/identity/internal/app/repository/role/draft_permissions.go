package role

import (
	"context"
	"errors"

	"grovehub/services/identity/internal/app/constants/permission"
)

var (
	ErrInvalidPermissions = errors.New("invalid or unavailable role permission set")
	ErrLiveRole          = errors.New("role is active or assigned; approved authority execution required")
)

type DraftPermissionUpdate struct {
	TenantID string
	ActorID  string
	RoleID   string
	Version  int64
	Codes    []permission.Code
}

// DraftPermissionWriter is intentionally separate from Repository and the role
// service graph. It is not an approval or authorization boundary and must not
// be exposed to a controller. It can only mutate inactive, unassigned roles.
type DraftPermissionWriter interface {
	ReplaceDraftPermissions(ctx context.Context, update DraftPermissionUpdate) (PermissionSet, error)
}

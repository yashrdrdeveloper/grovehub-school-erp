package rbacmigration

import (
	"context"
	"errors"

	"grovehub/services/identity/internal/app/constants/permission"
)

var ErrInvalidManifest = errors.New("invalid reviewed bootstrap manifest")
var ErrAdministrator = errors.New("initial administrator verification failed")

type RoleGrant struct {
	RoleID string `json:"roleId"`
	Codes []permission.Code `json:"permissionCodes"`
}

type Manifest struct {
	Plan Plan `json:"plan"`
	AdministratorID string `json:"administratorId"`
	Grants []RoleGrant `json:"grants"`
}

type BootstrapResult struct {
	AdministratorID string
	RolesActivated int
	PermissionsGranted int
}

// Bootstrap is a maintenance-only primitive. No HTTP or runtime service may
// call it. The operator must explicitly review the complete grant manifest.
type Bootstrap interface {
	Apply(ctx context.Context, manifest Manifest, operatorID string) (BootstrapResult, error)
}

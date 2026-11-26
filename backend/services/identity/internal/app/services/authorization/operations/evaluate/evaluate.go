// Package evaluate implements the pure school RBAC decision. Repository
// hydration and route/resource authorization must use this decision explicitly.
package evaluate

import (
	"context"

	coregateway "grovehub/libs/core/gateway"
	"grovehub/services/identity/internal/app/constants/feature"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/entity"
)

type RoleGrant struct {
	RoleID      string
	TenantID    string
	Active      bool
	Permissions []permission.Code
}

type PermissionState struct {
	Active  bool
	Feature feature.Code
}

type FeatureState struct {
	Active  bool
	Enabled bool
}

// Snapshot must be loaded from the authoritative Identity repositories, never
// from a request body or browser permission cache.
type Snapshot struct {
	AccountID   string
	TenantID    string
	Status      entity.AccountStatus
	Roles       []RoleGrant
	Permissions map[permission.Code]PermissionState
	Features    map[feature.Code]FeatureState
}

func Allowed(ctx context.Context, requested permission.Code, snapshot Snapshot) bool {
	principal, verified := coregateway.PrincipalFrom(ctx)
	if !verified || principal.TenantID == "" || principal.TenantID != snapshot.TenantID || principal.AccountID != snapshot.AccountID || snapshot.Status != entity.Active {
		return false
	}

	definition, known := permission.Find(requested)
	if !known {
		return false
	}
	state, registered := snapshot.Permissions[requested]
	if !registered || !state.Active {
		return false
	}
	if definition.Feature != state.Feature {
		return false
	}
	if state.Feature != "" {
		entitlement, enabled := snapshot.Features[state.Feature]
		if !feature.Known(state.Feature) || !enabled || !entitlement.Active || !entitlement.Enabled {
			return false
		}
	}

	for _, role := range snapshot.Roles {
		if !role.Active || role.RoleID == "" || role.TenantID != principal.TenantID {
			continue
		}
		for _, granted := range role.Permissions {
			if granted == requested {
				return true
			}
		}
	}
	return false
}

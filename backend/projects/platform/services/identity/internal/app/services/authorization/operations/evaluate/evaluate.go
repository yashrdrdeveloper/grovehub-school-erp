// Package evaluate implements the pure school RBAC decision. Repository
// hydration and route/resource authorization must use this decision explicitly.
package evaluate

import (
	"context"

	schoolgateway "grovehub/libs/schoolgateway"
	"grovehub/services/identity/internal/app/constants/feature"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/entity"
	authorizationrepo "grovehub/services/identity/internal/app/repository/authorization"
)

// Aliases preserve the evaluator's existing callers while the snapshot is owned by the repository.
type RoleGrant = authorizationrepo.RoleGrant
type PermissionState = authorizationrepo.PermissionState
type FeatureState = authorizationrepo.FeatureState
type Snapshot = authorizationrepo.Snapshot

func Allowed(ctx context.Context, requested permission.Code, snapshot Snapshot) bool {
	principal, verified := schoolgateway.PrincipalFrom(ctx)
	if !verified {
		return false
	}
	return AllowedFor(principal.AccountID, principal.TenantID, requested, snapshot)
}

// AllowedFor evaluates a database-loaded identity, not a browser credential.
// Only trusted callers may supply the subject; it never authenticates a session.
func AllowedFor(accountID string, tenantID string, requested permission.Code, snapshot Snapshot) bool {
	if accountID == "" || tenantID == "" || tenantID != snapshot.TenantID || accountID != snapshot.AccountID || snapshot.Status != entity.Active {
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
		if !role.Active || role.RoleID == "" || role.TenantID != tenantID {
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

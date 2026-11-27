package helpers

import (
	"grovehub/services/identity/internal/app/constants/feature"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/entity"
	authorizationrepo "grovehub/services/identity/internal/app/repository/authorization"
)

// SnapshotQuery reads a single account, its tenant-scoped assignments, registered
// permissions and the corresponding tenant entitlements in one SQL statement.
const SnapshotQuery = `SELECT
    a.account_id::text AS account_id,
    a.tenant_id AS tenant_id,
    a.status AS account_status,
    r.role_id::text AS role_id,
    r.tenant_id AS role_tenant_id,
    r.status AS role_status,
    p.code AS permission_code,
    p.active AS permission_active,
    f.code AS feature_code,
    f.active AS feature_active,
    tf.enabled AS feature_enabled
FROM identity.accounts a
LEFT JOIN identity.account_roles ar ON ar.tenant_id = a.tenant_id AND ar.account_id = a.account_id
LEFT JOIN identity.roles r ON r.tenant_id = ar.tenant_id AND r.role_id = ar.role_id
LEFT JOIN identity.role_permissions rp ON rp.tenant_id = r.tenant_id AND rp.role_id = r.role_id
LEFT JOIN identity.permissions p ON p.permission_id = rp.permission_id
LEFT JOIN identity.features f ON f.feature_id = p.feature_id
LEFT JOIN identity.tenant_features tf ON tf.tenant_id = a.tenant_id AND tf.feature_id = f.feature_id
WHERE a.tenant_id = ? AND a.account_id = ?
ORDER BY r.role_id, p.code`

// Null columns represent missing joins; absence is never treated as active.
type SnapshotRow struct {
	AccountID string
	TenantID string
	AccountStatus entity.AccountStatus
	RoleID *string
	RoleTenantID *string
	RoleStatus *string
	PermissionCode *permission.Code
	PermissionActive *bool
	FeatureCode *feature.Code
	FeatureActive *bool
	FeatureEnabled *bool
}

func HydrateSnapshot(rows []SnapshotRow) (authorizationrepo.Snapshot, error) {
	if len(rows) == 0 {
		return authorizationrepo.Snapshot{}, authorizationrepo.ErrNotFound
	}
	first := rows[0]
	snapshot := authorizationrepo.Snapshot{
		AccountID: first.AccountID,
		TenantID: first.TenantID,
		Status: first.AccountStatus,
		Permissions: make(map[permission.Code]authorizationrepo.PermissionState),
		Features: make(map[feature.Code]authorizationrepo.FeatureState),
	}
	roleIndexes := make(map[string]int)
	for _, row := range rows {
		if row.AccountID != first.AccountID || row.TenantID != first.TenantID || row.AccountStatus != first.AccountStatus {
			return authorizationrepo.Snapshot{}, authorizationrepo.ErrInvalidSnapshot
		}
		if err := appendGrant(&snapshot, roleIndexes, row); err != nil {
			return authorizationrepo.Snapshot{}, err
		}
	}
	return snapshot, nil
}

func appendGrant(snapshot *authorizationrepo.Snapshot, indexes map[string]int, row SnapshotRow) error {
	if row.RoleID == nil {
		return nil
	}
	if row.RoleTenantID == nil || *row.RoleTenantID != snapshot.TenantID || *row.RoleID == "" || row.RoleStatus == nil {
		return authorizationrepo.ErrInvalidSnapshot
	}
	index, exists := indexes[*row.RoleID]
	active := *row.RoleStatus == "ACTIVE"
	if !exists {
		index = len(snapshot.Roles)
		indexes[*row.RoleID] = index
		snapshot.Roles = append(snapshot.Roles, authorizationrepo.RoleGrant{RoleID: *row.RoleID, TenantID: *row.RoleTenantID, Active: active})
	}
	if snapshot.Roles[index].Active != active {
		return authorizationrepo.ErrInvalidSnapshot
	}
	if row.PermissionCode == nil {
		return nil
	}
	state := authorizationrepo.PermissionState{Active: isTrue(row.PermissionActive)}
	if row.FeatureCode != nil {
		state.Feature = *row.FeatureCode
		entitlement := authorizationrepo.FeatureState{Active: isTrue(row.FeatureActive), Enabled: isTrue(row.FeatureEnabled)}
		if previous, exists := snapshot.Features[state.Feature]; exists && previous != entitlement {
			return authorizationrepo.ErrInvalidSnapshot
		}
		snapshot.Features[state.Feature] = entitlement
	}
	if previous, exists := snapshot.Permissions[*row.PermissionCode]; exists && previous != state {
		return authorizationrepo.ErrInvalidSnapshot
	}
	snapshot.Permissions[*row.PermissionCode] = state
	snapshot.Roles[index].Permissions = append(snapshot.Roles[index].Permissions, *row.PermissionCode)
	return nil
}

func isTrue(value *bool) bool {
	return value != nil && *value
}

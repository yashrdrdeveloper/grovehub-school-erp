package helpers

import (
	"slices"

	"grovehub/services/identity/internal/app/constants/permission"
	dto "grovehub/services/identity/internal/app/dto/capability"
	repo "grovehub/services/identity/internal/app/repository/capability"
	"grovehub/services/identity/internal/app/services/security/operations/routepolicy"
)

// Only advertise routes that actually exist and are not blocked. These are
// UI hints, not resource-specific authorization or approval eligibility.
func Actions(effective []permission.Code) map[string]bool {
	patterns := map[string]string{
		"auditSearch": "POST /audit/protected/search",
		"accountSearch": "POST /identity/protected/accounts/search",
		"accountView": "GET /identity/protected/accounts/{accountId}",
		"accountCreate": "POST /identity/protected/accounts",
		"accountUpdate": "PATCH /identity/protected/accounts/{accountId}",
		"accountResendInvitation": "POST /identity/protected/accounts/{accountId}/resend-invitation",
		"roleSearch": "POST /identity/protected/roles/search",
		"roleView": "GET /identity/protected/roles/{roleId}",
		"roleCreate": "POST /identity/protected/roles",
		"roleUpdate": "PATCH /identity/protected/roles/{roleId}",
		"rolePermissionsView": "GET /identity/protected/roles/{roleId}/permissions",
		"rolePermissionsSubmit": "POST /identity/protected/approvals/permissions",
		"accountRolesView": "GET /identity/protected/account-roles/{accountId}",
		"accountRolesSubmit": "POST /identity/protected/approvals/assignments",
		"permissionCatalog": "GET /identity/protected/capabilities/permissions",
	}
	actions := make(map[string]bool, len(patterns)+5)
	for name, pattern := range patterns {
		rule, known := routepolicy.Lookup(pattern)
		actions[name] = known && !rule.BlockedUntilCutover && rule.Permission != "" && slices.Contains(effective, rule.Permission)
	}
	actions["accountSuspendRequest"] = slices.Contains(effective, permission.AccountSuspend)
	actions["accountActivateRequest"] = slices.Contains(effective, permission.AccountActivate)
	actions["accountDeactivateRequest"] = slices.Contains(effective, permission.AccountDeactivate)
	actions["roleActivateRequest"] = slices.Contains(effective, permission.RoleUpdate)
	actions["roleDeactivateRequest"] = slices.Contains(effective, permission.RoleDeactivate)
	return actions
}

func Catalog(rows []repo.Entry, effective []permission.Code) dto.Catalog {
	result := dto.Catalog{Permissions: make([]dto.Permission, 0, len(rows))}
	canManage := slices.Contains(effective, permission.RolePermissionManage)
	for _, row := range rows {
		definition, known := permission.Find(row.Code)
		if !known || !permission.SchoolDelegable(row.Code) || definition.Feature != row.Feature {
			continue
		}
		result.Permissions = append(result.Permissions, dto.Permission{Code: row.Code, Name: definition.Name, Feature: definition.Feature, CanDelegate: canManage && slices.Contains(effective, row.Code)})
	}
	return result
}

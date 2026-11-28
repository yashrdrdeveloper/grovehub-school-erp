package manifest

import (
	"slices"

	"github.com/google/uuid"

	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/services/rbacmigration"
)

func Validate(request rbacmigration.Manifest) error {
	id, err := uuid.Parse(request.AdministratorID)
	if err != nil || id == uuid.Nil || id.String() != request.AdministratorID || len(request.Grants) == 0 {
		return rbacmigration.ErrInvalidManifest
	}
	roles := make(map[string]bool)
	for _, role := range request.Plan.Roles {
		roles[role.RoleID] = true
	}
	owners := make(map[string][]string)
	for _, assignment := range request.Plan.Assignments {
		for _, roleID := range append([]string{assignment.PrimaryRoleID}, assignment.SecondaryRoleIDs...) {
			owners[roleID] = append(owners[roleID], assignment.AccountID)
		}
	}
	seen := make(map[string]bool)
	recovery := permission.AdministratorRecoveryCapabilities()
	held := make(map[permission.Code]bool)
	for _, grant := range request.Grants {
		if !roles[grant.RoleID] || seen[grant.RoleID] || len(grant.Codes) == 0 || len(grant.Codes) > len(permission.Catalog()) {
			return rbacmigration.ErrInvalidManifest
		}
		seen[grant.RoleID] = true
		codes := make(map[permission.Code]bool)
		for _, code := range grant.Codes {
			if !permission.SchoolDelegable(code) || codes[code] {
				return rbacmigration.ErrInvalidManifest
			}
			codes[code] = true
			if slices.Contains(owners[grant.RoleID], request.AdministratorID) {
				held[code] = true
			}
			// Do not accidentally grant recovery authority to other legacy
			// role holders (including suspended or pending accounts).
			if slices.Contains(recovery, code) && (len(owners[grant.RoleID]) != 1 || owners[grant.RoleID][0] != request.AdministratorID) {
				return rbacmigration.ErrAdministrator
			}
		}
	}
	for _, code := range recovery {
		if !held[code] {
			return rbacmigration.ErrAdministrator
		}
	}
	return nil
}

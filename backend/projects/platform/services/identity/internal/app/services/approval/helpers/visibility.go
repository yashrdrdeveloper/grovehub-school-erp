package helpers

import (
	"grovehub/services/identity/internal/app/entity"
	repo "grovehub/services/identity/internal/app/repository/approval"
	authrepo "grovehub/services/identity/internal/app/repository/authorization"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/authorization/operations/evaluate"
)

func EligibleGrants(accountID string, tenantID string, snapshot authrepo.Snapshot) ([]repo.EligibleGrant, error) {
	if accountID == "" || tenantID == "" || snapshot.AccountID != accountID || snapshot.TenantID != tenantID || snapshot.Status != entity.Active {
		return nil, authorization.ErrForbidden
	}
	grants := make([]repo.EligibleGrant, 0)
	for _, role := range snapshot.Roles {
		if !role.Active || role.RoleID == "" || role.TenantID != tenantID {
			continue
		}
		scoped := snapshot
		scoped.Roles = []authrepo.RoleGrant{role}
		for _, code := range role.Permissions {
			if evaluate.AllowedFor(accountID, tenantID, code, scoped) {
				grants = append(grants, repo.EligibleGrant{RoleID: role.RoleID, Permission: code})
			}
		}
	}
	return grants, nil
}

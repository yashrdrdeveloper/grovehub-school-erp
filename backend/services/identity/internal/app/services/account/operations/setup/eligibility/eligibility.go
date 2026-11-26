package eligibility

import (
	"context"
	"slices"

	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/repository/accountrole"
	authrepo "grovehub/services/identity/internal/app/repository/authorization"
	"grovehub/services/identity/internal/app/services/authorization/operations/evaluate"
)

// Require validates the normalized role set in a tenant authority transaction.
// The prospective status is used only for pure evaluation, never as a principal.
func Require(ctx context.Context, account entity.Account, assignments accountrole.Repository, snapshots authrepo.Repository) error {
	if account.TenantID == nil || *account.TenantID == "" || account.Status != entity.PendingSetup || assignments == nil || snapshots == nil {
		return accountrole.ErrIneligible
	}
	tenant := *account.TenantID
	assignment, err := assignments.Find(ctx, tenant, account.AccountID)
	if err != nil {
		return err
	}
	if assignment.AccountID != account.AccountID || assignment.Version != account.Version || assignment.Status != entity.PendingSetup || assignment.PrimaryRoleID == "" {
		return accountrole.ErrIneligible
	}
	ids := append([]string{assignment.PrimaryRoleID}, assignment.SecondaryRoleIDs...)
	snapshot, err := snapshots.LoadSnapshot(ctx, tenant, account.AccountID)
	if err != nil {
		return err
	}
	if snapshot.AccountID != account.AccountID || snapshot.TenantID != tenant || snapshot.Status != entity.PendingSetup || len(snapshot.Roles) != len(ids) {
		return accountrole.ErrIneligible
	}
	snapshot.Status = entity.Active
	seen := make(map[string]bool, len(ids))
	for _, role := range snapshot.Roles {
		if !role.Active || role.TenantID != tenant || seen[role.RoleID] || !slices.Contains(ids, role.RoleID) || len(role.Permissions) == 0 {
			return accountrole.ErrIneligible
		}
		seen[role.RoleID] = true
		scoped := snapshot
		scoped.Roles = []authrepo.RoleGrant{role}
		for _, code := range role.Permissions {
			if !permission.SchoolDelegable(code) || !evaluate.AllowedFor(account.AccountID, tenant, code, scoped) {
				return accountrole.ErrIneligible
			}
		}
	}
	return nil
}

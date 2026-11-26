package plan

import (
	"encoding/json"
	"slices"
	"strings"

	"github.com/google/uuid"

	"grovehub/services/identity/internal/app/constants/persona"
	"grovehub/services/identity/internal/app/services/rbacmigration"
)

// A fixed, domain-specific namespace makes the mapping stable across retries.
// The JSON tuple avoids ambiguous tenant/code concatenation.
var namespace = uuid.MustParse("28a81cc0-273e-4e49-88ee-d88789436580")

func RoleID(tenantID string, code persona.Code) string {
	key, _ := json.Marshal([]string{tenantID, string(code)})
	return uuid.NewSHA1(namespace, key).String()
}

// Build is deliberately strict: an existing normalized assignment, role ID or
// name collision blocks migration rather than silently adopting existing data.
func Build(inventory rbacmigration.Inventory) (rbacmigration.Plan, error) {
	tenant := inventory.TenantID
	if strings.TrimSpace(tenant) == "" || strings.TrimSpace(tenant) != tenant {
		return rbacmigration.Plan{}, rbacmigration.ErrScopeRequired
	}
	if inventory.NormalizedAssignments != 0 {
		return rbacmigration.Plan{}, rbacmigration.ErrConflict
	}
	result := rbacmigration.Plan{TenantID: tenant, Roles: []rbacmigration.RolePlan{}, Assignments: []rbacmigration.AssignmentPlan{}, UnassignedAccountIDs: []string{}}
	codes := make(map[string]string)
	for _, definition := range persona.Catalog() {
		id := RoleID(tenant, definition.Code)
		for _, existing := range inventory.Roles {
			if existing.RoleID == id || existing.NormalizedName == strings.ToLower(definition.Name) {
				return rbacmigration.Plan{}, rbacmigration.ErrConflict
			}
		}
		codes[string(definition.Code)] = id
		result.Roles = append(result.Roles, rbacmigration.RolePlan{RoleID: id, Name: definition.Name, Persona: definition.Code})
	}
	seen := make(map[string]bool)
	for _, account := range inventory.Accounts {
		id, err := uuid.Parse(account.AccountID)
		if err != nil || id == uuid.Nil || id.String() != account.AccountID || account.Version < 1 || seen[account.AccountID] {
			return rbacmigration.Plan{}, rbacmigration.ErrInvalidInventory
		}
		seen[account.AccountID] = true
		if account.Primary == nil {
			if len(account.Secondary) != 0 {
				return rbacmigration.Plan{}, rbacmigration.ErrInvalidInventory
			}
			result.UnassignedAccountIDs = append(result.UnassignedAccountIDs, account.AccountID)
			continue
		}
		primary, known := codes[*account.Primary]
		if !known {
			return rbacmigration.Plan{}, rbacmigration.ErrInvalidInventory
		}
		assignment := rbacmigration.AssignmentPlan{AccountID: account.AccountID, Version: account.Version, PrimaryRoleID: primary, SecondaryRoleIDs: []string{}}
		assigned := map[string]bool{primary: true}
		for _, code := range account.Secondary {
			id, known := codes[code]
			if !known || assigned[id] {
				return rbacmigration.Plan{}, rbacmigration.ErrInvalidInventory
			}
			assigned[id] = true
			assignment.SecondaryRoleIDs = append(assignment.SecondaryRoleIDs, id)
		}
		slices.Sort(assignment.SecondaryRoleIDs)
		result.Assignments = append(result.Assignments, assignment)
	}
	slices.SortFunc(result.Assignments, func(a, b rbacmigration.AssignmentPlan) int {
		return strings.Compare(a.AccountID, b.AccountID)
	})
	slices.Sort(result.UnassignedAccountIDs)
	return result, nil
}

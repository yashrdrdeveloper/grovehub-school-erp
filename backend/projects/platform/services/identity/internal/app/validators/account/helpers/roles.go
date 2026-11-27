package helpers

import (
	"errors"
	"strings"

	roles "grovehub/services/identity/internal/app/constants/role"
	accountdto "grovehub/services/identity/internal/app/dto/account"
)

func NormalizeRoleAssignment(input accountdto.RoleAssignment) (accountdto.RoleAssignment, error) {
	primary := strings.TrimSpace(input.PrimaryRoleID)
	if !roles.IsAssignable(primary) {
		return accountdto.RoleAssignment{}, errors.New("invalid primary role")
	}
	catalog := roles.Catalog()
	if len(input.SecondaryRoleIDs) >= len(catalog) {
		return accountdto.RoleAssignment{}, errors.New("invalid secondary roles")
	}
	seen := map[string]bool{primary: true}
	secondary := make([]string, 0, len(input.SecondaryRoleIDs))
	for _, value := range input.SecondaryRoleIDs {
		id := strings.TrimSpace(value)
		if !roles.IsAssignable(id) || seen[id] {
			return accountdto.RoleAssignment{}, errors.New("invalid or duplicate secondary role")
		}
		seen[id] = true
		secondary = append(secondary, id)
	}
	// The secondary role set is order-independent, including in idempotency hashes.
	ordered := make([]string, 0, len(secondary))
	for _, item := range catalog {
		if item.ID != primary && seen[item.ID] {
			ordered = append(ordered, item.ID)
		}
	}
	return accountdto.RoleAssignment{PrimaryRoleID: primary, SecondaryRoleIDs: ordered}, nil
}

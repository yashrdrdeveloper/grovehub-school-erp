package helpers

import (
	"slices"
	"strings"

	"github.com/google/uuid"

	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/repository/accountrole"
)

func NormalizeReplacement(request accountrole.Replacement) (accountrole.Replacement, error) {
	if strings.TrimSpace(request.TenantID) == "" || strings.TrimSpace(request.TenantID) != request.TenantID {
		return accountrole.Replacement{}, accountrole.ErrScopeRequired
	}
	if request.Version < 1 || len(request.SecondaryRoleIDs) > 20 || len(request.RoleVersions) > 42 {
		return accountrole.Replacement{}, accountrole.ErrInvalidAssignment
	}
	for _, value := range []string{request.ActorID, request.AccountID, request.PrimaryRoleID} {
		if !canonicalID(value) {
			return accountrole.Replacement{}, accountrole.ErrInvalidAssignment
		}
	}
	seen := map[string]bool{request.PrimaryRoleID: true}
	for _, id := range request.SecondaryRoleIDs {
		if !canonicalID(id) || seen[id] {
			return accountrole.Replacement{}, accountrole.ErrInvalidAssignment
		}
		seen[id] = true
	}
	versions := make(map[string]bool)
	for _, item := range request.RoleVersions {
		if !canonicalID(item.RoleID) || item.Version < 1 || versions[item.RoleID] {
			return accountrole.Replacement{}, accountrole.ErrInvalidAssignment
		}
		versions[item.RoleID] = true
	}
	request.SecondaryRoleIDs = slices.Clone(request.SecondaryRoleIDs)
	slices.Sort(request.SecondaryRoleIDs)
	return request, nil
}

func canonicalID(value string) bool {
	id, err := uuid.Parse(value)
	return err == nil && id != uuid.Nil && id.String() == value
}

func RoleIDs(current accountrole.Assignment, request accountrole.Replacement) []string {
	ids := append([]string{}, current.SecondaryRoleIDs...)
	if current.PrimaryRoleID != "" {
		ids = append(ids, current.PrimaryRoleID)
	}
	ids = append(ids, request.PrimaryRoleID)
	ids = append(ids, request.SecondaryRoleIDs...)
	slices.Sort(ids)
	return slices.Compact(ids)
}

func ValidateCandidates(current accountrole.Assignment, request accountrole.Replacement, candidates []accountrole.Candidate) error {
	if current.Version != request.Version {
		return accountrole.ErrVersionConflict
	}
	if current.Status != entity.Active && current.Status != entity.PendingSetup {
		return accountrole.ErrIneligible
	}
	if current.PrimaryRoleID == "" && len(current.SecondaryRoleIDs) > 0 {
		return accountrole.ErrInvalidAssignment
	}
	ids := RoleIDs(current, request)
	if len(ids) != len(candidates) || len(ids) != len(request.RoleVersions) {
		return accountrole.ErrVersionConflict
	}
	versions := make(map[string]int64)
	for _, item := range request.RoleVersions {
		versions[item.RoleID] = item.Version
	}
	seen := make(map[string]bool)
	for _, candidate := range candidates {
		if seen[candidate.RoleID] || !slices.Contains(ids, candidate.RoleID) || versions[candidate.RoleID] != candidate.Version {
			return accountrole.ErrVersionConflict
		}
		seen[candidate.RoleID] = true
		if candidate.RoleID == request.PrimaryRoleID || slices.Contains(request.SecondaryRoleIDs, candidate.RoleID) {
			if candidate.Status != entity.RoleActive || len(candidate.Codes) == 0 {
				return accountrole.ErrIneligible
			}
		}
	}
	return nil
}

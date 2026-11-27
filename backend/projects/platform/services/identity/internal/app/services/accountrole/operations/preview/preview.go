package preview

import (
	"slices"

	"grovehub/services/identity/internal/app/constants/permission"
	dto "grovehub/services/identity/internal/app/dto/accountrole"
	"grovehub/services/identity/internal/app/entity"
	repo "grovehub/services/identity/internal/app/repository/accountrole"
	"grovehub/services/identity/internal/app/services/authorization"
)

// RoleIDs returns the union of old and proposed role IDs for one batch lookup.
func RoleIDs(current repo.Assignment, request dto.Replace) []string {
	ids := append([]string{}, current.SecondaryRoleIDs...)
	if current.PrimaryRoleID != "" {
		ids = append(ids, current.PrimaryRoleID)
	}
	ids = append(ids, request.PrimaryRoleID)
	ids = append(ids, request.SecondaryRoleIDs...)
	slices.Sort(ids)
	return slices.Compact(ids)
}

func Evaluate(current repo.Assignment, request dto.Replace, candidates []repo.Candidate, effective []permission.Code) (dto.Preview, error) {
	if !slices.Contains(effective, permission.AccountRoleAssign) {
		return dto.Preview{}, authorization.ErrForbidden
	}
	if current.Version != request.Version {
		return dto.Preview{}, repo.ErrVersionConflict
	}
	if current.Status != entity.Active && current.Status != entity.PendingSetup {
		return dto.Preview{}, repo.ErrIneligible
	}
	if current.PrimaryRoleID == "" && len(current.SecondaryRoleIDs) > 0 {
		return dto.Preview{}, repo.ErrInvalidAssignment
	}
	expected := RoleIDs(current, request)
	index := make(map[string]repo.Candidate, len(candidates))
	for _, candidate := range candidates {
		if _, exists := index[candidate.RoleID]; exists || candidate.Version < 1 {
			return dto.Preview{}, repo.ErrInvalidAssignment
		}
		index[candidate.RoleID] = candidate
	}
	if len(index) != len(expected) {
		return dto.Preview{}, repo.ErrNotFound
	}
	versions := make([]dto.RoleVersion, 0, len(expected))
	for _, id := range expected {
		candidate, exists := index[id]
		if !exists {
			return dto.Preview{}, repo.ErrNotFound
		}
		if (id == request.PrimaryRoleID || slices.Contains(request.SecondaryRoleIDs, id)) && (candidate.Status != entity.RoleActive || len(candidate.Codes) == 0) {
			return dto.Preview{}, repo.ErrIneligible
		}
		for _, code := range candidate.Codes {
			if !permission.SchoolDelegable(code) || !slices.Contains(effective, code) {
				return dto.Preview{}, authorization.ErrForbidden
			}
		}
		versions = append(versions, dto.RoleVersion{RoleID: id, Version: candidate.Version})
	}
	return dto.Preview{
		AccountID: current.AccountID, Version: current.Version,
		PrimaryRoleID: request.PrimaryRoleID, SecondaryRoleIDs: request.SecondaryRoleIDs,
		RoleVersions: versions,
		ApprovalRequired: current.PrimaryRoleID != request.PrimaryRoleID || !slices.Equal(current.SecondaryRoleIDs, request.SecondaryRoleIDs),
	}, nil
}

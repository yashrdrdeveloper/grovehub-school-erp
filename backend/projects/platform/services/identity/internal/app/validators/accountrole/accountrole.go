package accountrole

import (
	"slices"

	"github.com/google/uuid"

	dto "grovehub/services/identity/internal/app/dto/accountrole"
	repo "grovehub/services/identity/internal/app/repository/accountrole"
)

const MaxSecondaryRoles = 20

func NormalizeID(value string) (string, error) {
	id, err := uuid.Parse(value)
	if err != nil || id == uuid.Nil {
		return "", repo.ErrInvalidAssignment
	}
	return id.String(), nil
}

func NormalizeReplace(request dto.Replace) (dto.Replace, error) {
	if request.Version < 1 || len(request.SecondaryRoleIDs) > MaxSecondaryRoles {
		return dto.Replace{}, repo.ErrInvalidAssignment
	}
	primary, err := NormalizeID(request.PrimaryRoleID)
	if err != nil {
		return dto.Replace{}, err
	}
	seen := map[string]bool{primary: true}
	secondary := make([]string, 0, len(request.SecondaryRoleIDs))
	for _, value := range request.SecondaryRoleIDs {
		id, err := NormalizeID(value)
		if err != nil || seen[id] {
			return dto.Replace{}, repo.ErrInvalidAssignment
		}
		seen[id] = true
		secondary = append(secondary, id)
	}
	slices.Sort(secondary)
	request.PrimaryRoleID = primary
	request.SecondaryRoleIDs = secondary
	return request, nil
}

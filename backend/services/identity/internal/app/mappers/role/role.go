package role

import (
	"grovehub/services/identity/internal/app/constants/persona"
	roledto "grovehub/services/identity/internal/app/dto/role"
	rolerepo "grovehub/services/identity/internal/app/repository/role"
)

func ToDetail(detail rolerepo.Detail) roledto.Detail {
	return roledto.Detail{
		RoleID:      detail.Role.RoleID,
		Name:        detail.Role.Name,
		Description: detail.Role.Description,
		Status:      string(detail.Role.Status),
		Personas:    append([]persona.Code{}, detail.Personas...),
		Version:     detail.Role.Version,
		CreatedAt:   detail.Role.CreatedAt,
		UpdatedAt:   detail.Role.UpdatedAt,
	}
}

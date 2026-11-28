package helpers

import (
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"grovehub/services/identity/internal/app/constants/feature"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/constants/persona"
	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/repository/role"
	"grovehub/services/identity/internal/app/repository/role/constants"
)

func NormalizeTransition(request role.Transition) (role.Transition, error) {
	if strings.TrimSpace(request.TenantID) == "" || strings.TrimSpace(request.TenantID) != request.TenantID {
		return role.Transition{}, role.ErrScopeRequired
	}
	actor, err := uuid.Parse(request.ActorID)
	if err != nil || actor == uuid.Nil {
		return role.Transition{}, role.ErrInvalidTransition
	}
	id, err := uuid.Parse(request.RoleID)
	if err != nil || id == uuid.Nil {
		return role.Transition{}, role.ErrNotFound
	}
	if request.Version < 1 {
		return role.Transition{}, role.ErrVersionConflict
	}
	if request.Status != entity.RoleActive && request.Status != entity.RoleInactive {
		return role.Transition{}, role.ErrInvalidTransition
	}
	request.ActorID = actor.String()
	request.RoleID = id.String()
	return request, nil
}

func CheckActivation(db *gorm.DB, detail role.Detail, codes []permission.Code) error {
	if len(detail.Personas) == 0 || len(codes) == 0 {
		return role.ErrNotReady
	}
	for _, code := range detail.Personas {
		if !persona.Known(code) {
			return role.ErrNotReady
		}
	}
	required := make(map[feature.Code]bool)
	for _, code := range codes {
		definition, known := permission.Find(code)
		if !known || !permission.SchoolDelegable(code) {
			return role.ErrInvalidPermissions
		}
		if definition.Feature != "" {
			required[definition.Feature] = true
		}
	}
	if _, err := ResolvePermissionIDs(db, codes); err != nil {
		return err
	}
	for _, definition := range feature.Catalog() {
		if !required[definition.Code] {
			continue
		}
		var count int64
		if err := db.Raw(constants.EnabledFeature, detail.Role.TenantID, definition.Code).Scan(&count).Error; err != nil {
			return err
		}
		if count != 1 {
			return role.ErrNotReady
		}
		delete(required, definition.Code)
	}
	if len(required) != 0 {
		return role.ErrNotReady
	}
	return nil
}

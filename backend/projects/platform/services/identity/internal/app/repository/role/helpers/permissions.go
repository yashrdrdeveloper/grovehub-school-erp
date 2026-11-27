package helpers

import (
	"slices"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/repository/role"
	"grovehub/services/identity/internal/app/repository/role/constants"
)

func NormalizePermissionUpdate(update role.DraftPermissionUpdate) (role.DraftPermissionUpdate, error) {
	if strings.TrimSpace(update.TenantID) == "" || strings.TrimSpace(update.TenantID) != update.TenantID {
		return role.DraftPermissionUpdate{}, role.ErrScopeRequired
	}
	actorID, err := uuid.Parse(update.ActorID)
	if err != nil || actorID == uuid.Nil {
		return role.DraftPermissionUpdate{}, role.ErrInvalidPermissions
	}
	roleID, err := uuid.Parse(update.RoleID)
	if err != nil || roleID == uuid.Nil {
		return role.DraftPermissionUpdate{}, role.ErrNotFound
	}
	if update.Version < 1 {
		return role.DraftPermissionUpdate{}, role.ErrVersionConflict
	}
	if update.Codes == nil || len(update.Codes) > len(permission.Catalog()) {
		return role.DraftPermissionUpdate{}, role.ErrInvalidPermissions
	}
	seen := make(map[permission.Code]bool, len(update.Codes))
	for _, code := range update.Codes {
		if !permission.SchoolDelegable(code) || seen[code] {
			return role.DraftPermissionUpdate{}, role.ErrInvalidPermissions
		}
		seen[code] = true
	}
	update.ActorID = actorID.String()
	update.RoleID = roleID.String()
	update.Codes = slices.Clone(update.Codes)
	slices.Sort(update.Codes)
	return update, nil
}

type registeredPermission struct {
	ID          string
	Module      string
	FeatureCode string
}

func ResolvePermissionIDs(db *gorm.DB, codes []permission.Code) ([]string, error) {
	ids := make([]string, 0, len(codes))
	for _, code := range codes {
		definition, known := permission.Find(code)
		if !known {
			return nil, role.ErrInvalidPermissions
		}
		var stored registeredPermission
		result := db.Raw(constants.FindActivePermission, code).Scan(&stored)
		if result.Error != nil {
			return nil, result.Error
		}
		if result.RowsAffected != 1 || stored.ID == "" || stored.Module != "IDENTITY" || stored.FeatureCode != string(definition.Feature) {
			return nil, role.ErrInvalidPermissions
		}
		ids = append(ids, stored.ID)
	}
	return ids, nil
}

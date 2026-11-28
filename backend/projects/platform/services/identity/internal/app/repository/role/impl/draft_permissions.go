package impl

import (
	"context"
	"slices"

	"gorm.io/gorm"

	coredatabase "grovehub/libs/core/database"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/repository/authority"
	authorityimpl "grovehub/services/identity/internal/app/repository/authority/impl"
	"grovehub/services/identity/internal/app/repository/role"
	"grovehub/services/identity/internal/app/repository/role/constants"
	"grovehub/services/identity/internal/app/repository/role/helpers"
)

type DraftWriter struct {
	db     *gorm.DB
	locker authority.Locker
	reads  role.Repository
}

var _ role.DraftPermissionWriter = (*DraftWriter)(nil)

func NewDraftWriter(db *gorm.DB) role.DraftPermissionWriter {
	return &DraftWriter{db: db, locker: authorityimpl.New(db), reads: New(db)}
}

type lockedRole struct {
	Status  entity.RoleStatus
	Version int64
}

func (writer *DraftWriter) ReplaceDraftPermissions(ctx context.Context, update role.DraftPermissionUpdate) (role.PermissionSet, error) {
	update, err := helpers.NormalizePermissionUpdate(update)
	if err != nil {
		return role.PermissionSet{}, err
	}
	var detail role.PermissionSet
	err = writer.locker.Run(ctx, update.TenantID, func(tx context.Context) error {
		db := coredatabase.TransactionDB(tx, writer.db)
		if err := checkDraft(db, update); err != nil {
			return err
		}
		current, err := writer.reads.Permissions(tx, update.TenantID, update.RoleID)
		if err != nil {
			return err
		}
		for _, code := range current.Codes {
			if !permission.SchoolDelegable(code) {
				return role.ErrInvalidPermissions
			}
		}
		ids, err := helpers.ResolvePermissionIDs(db, update.Codes)
		if err != nil {
			return err
		}
		if slices.Equal(current.Codes, update.Codes) {
			detail = current
			return nil
		}
		if err := replacePermissions(db, update, ids); err != nil {
			return err
		}
		detail, err = writer.reads.Permissions(tx, update.TenantID, update.RoleID)
		return err
	})
	if err != nil {
		return role.PermissionSet{}, err
	}
	return detail, nil
}

func checkDraft(db *gorm.DB, update role.DraftPermissionUpdate) error {
	var locked lockedRole
	result := db.Raw(constants.LockRole, update.TenantID, update.RoleID).Scan(&locked)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return role.ErrNotFound
	}
	if locked.Version != update.Version {
		return role.ErrVersionConflict
	}
	if locked.Status != entity.RoleInactive {
		return role.ErrLiveRole
	}
	var count int64
	if err := db.Raw(constants.CountAssignments, update.TenantID, update.RoleID).Scan(&count).Error; err != nil {
		return err
	}
	if count != 0 {
		return role.ErrLiveRole
	}
	return nil
}

func replacePermissions(db *gorm.DB, update role.DraftPermissionUpdate, ids []string) error {
	if err := db.Exec(constants.DeletePermissions, update.TenantID, update.RoleID).Error; err != nil {
		return err
	}
	for _, id := range ids {
		if err := db.Exec(constants.InsertPermission, update.TenantID, update.RoleID, id, update.ActorID).Error; err != nil {
			return err
		}
	}
	result := db.Exec(constants.AdvancePermissionVersion, update.ActorID, update.TenantID, update.RoleID, update.Version)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return role.ErrVersionConflict
	}
	return nil
}

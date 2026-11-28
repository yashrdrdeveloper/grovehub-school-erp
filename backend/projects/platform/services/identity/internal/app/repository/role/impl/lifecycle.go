package impl

import (
	"context"

	"gorm.io/gorm"

	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/repository/authority"
	authorityimpl "grovehub/services/identity/internal/app/repository/authority/impl"
	"grovehub/services/identity/internal/app/repository/role"
	"grovehub/services/identity/internal/app/repository/role/constants"
	"grovehub/services/identity/internal/app/repository/role/helpers"
)

type LifecycleWriter struct {
	db     *gorm.DB
	reads  role.Repository
	locker authority.Locker
}

var _ role.LifecycleWriter = (*LifecycleWriter)(nil)

func NewLifecycleWriter(db *gorm.DB) role.LifecycleWriter {
	return &LifecycleWriter{db: db, reads: New(db), locker: authorityimpl.New(db)}
}

func (writer *LifecycleWriter) Transition(ctx context.Context, request role.Transition) (role.Detail, error) {
	request, err := helpers.NormalizeTransition(request)
	if err != nil {
		return role.Detail{}, err
	}
	var detail role.Detail
	err = writer.locker.Run(ctx, request.TenantID, func(tx context.Context) error {
		db := coredatabase.TransactionDB(tx, writer.db)
		var locked lockedRole
		result := db.Raw(constants.LockRole, request.TenantID, request.RoleID).Scan(&locked)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return role.ErrNotFound
		}
		if locked.Version != request.Version {
			return role.ErrVersionConflict
		}
		if locked.Status == request.Status || (locked.Status != entity.RoleActive && locked.Status != entity.RoleInactive) {
			return role.ErrInvalidTransition
		}
		var assigned int64
		if err := db.Raw(constants.CountAssignments, request.TenantID, request.RoleID).Scan(&assigned).Error; err != nil {
			return err
		}
		if assigned != 0 {
			return role.ErrAssigned
		}
		if request.Status == entity.RoleActive {
			current, err := writer.reads.FindByID(tx, request.TenantID, request.RoleID)
			if err != nil {
				return err
			}
			grants, err := writer.reads.Permissions(tx, request.TenantID, request.RoleID)
			if err != nil {
				return err
			}
			if err := helpers.CheckActivation(db, current, grants.Codes); err != nil {
				return err
			}
		}
		result = db.Exec(constants.TransitionStatus, request.Status, request.ActorID, request.TenantID, request.RoleID, request.Version, locked.Status)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return role.ErrVersionConflict
		}
		var err error
		detail, err = writer.reads.FindByID(tx, request.TenantID, request.RoleID)
		return err
	})
	if err != nil {
		return role.Detail{}, err
	}
	return detail, nil
}

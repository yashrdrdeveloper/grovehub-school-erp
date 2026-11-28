package impl

import (
	"context"
	"slices"

	"gorm.io/gorm"

	coredatabase "grovehub/libs/core/database"
	"grovehub/services/identity/internal/app/repository/accountrole"
	"grovehub/services/identity/internal/app/repository/accountrole/constants"
	"grovehub/services/identity/internal/app/repository/accountrole/helpers"
	"grovehub/services/identity/internal/app/repository/authority"
	authorityimpl "grovehub/services/identity/internal/app/repository/authority/impl"
)

type Writer struct {
	db     *gorm.DB
	locker authority.Locker
	reads  accountrole.Repository
}

var _ accountrole.Writer = (*Writer)(nil)

func NewWriter(db *gorm.DB) accountrole.Writer {
	return &Writer{db: db, locker: authorityimpl.New(db), reads: New(db)}
}

func (writer *Writer) Replace(ctx context.Context, request accountrole.Replacement) (accountrole.Assignment, error) {
	request, err := helpers.NormalizeReplacement(request)
	if err != nil {
		return accountrole.Assignment{}, err
	}
	var assignment accountrole.Assignment
	err = writer.locker.Run(ctx, request.TenantID, func(tx context.Context) error {
		db := coredatabase.TransactionDB(tx, writer.db)
		var locked string
		result := db.Raw(constants.LockAccount, request.TenantID, request.AccountID).Scan(&locked)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return accountrole.ErrNotFound
		}
		current, err := writer.reads.Find(tx, request.TenantID, request.AccountID)
		if err != nil {
			return err
		}
		ids := helpers.RoleIDs(current, request)
		var lockedIDs []string
		if err := db.Raw(constants.LockCandidates, request.TenantID, ids).Scan(&lockedIDs).Error; err != nil {
			return err
		}
		if len(lockedIDs) != len(ids) {
			return accountrole.ErrNotFound
		}
		candidates, err := writer.reads.Candidates(tx, request.TenantID, ids)
		if err != nil {
			return err
		}
		if err := helpers.ValidateCandidates(current, request, candidates); err != nil {
			return err
		}
		if current.PrimaryRoleID == request.PrimaryRoleID && slices.Equal(current.SecondaryRoleIDs, request.SecondaryRoleIDs) {
			assignment = current
			return nil
		}
		if err := replaceAssignments(db, request); err != nil {
			return err
		}
		assignment, err = writer.reads.Find(tx, request.TenantID, request.AccountID)
		return err
	})
	if err != nil {
		return accountrole.Assignment{}, err
	}
	return assignment, nil
}

func replaceAssignments(db *gorm.DB, request accountrole.Replacement) error {
	if err := db.Exec(constants.DeleteAssignments, request.TenantID, request.AccountID).Error; err != nil {
		return err
	}
	if err := db.Exec(constants.InsertAssignment, request.TenantID, request.AccountID, request.PrimaryRoleID, "PRIMARY", request.ActorID).Error; err != nil {
		return err
	}
	for _, id := range request.SecondaryRoleIDs {
		if err := db.Exec(constants.InsertAssignment, request.TenantID, request.AccountID, id, "SECONDARY", request.ActorID).Error; err != nil {
			return err
		}
	}
	result := db.Exec(constants.AdvanceAccountVersion, request.TenantID, request.AccountID, request.Version)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return accountrole.ErrVersionConflict
	}
	return nil
}

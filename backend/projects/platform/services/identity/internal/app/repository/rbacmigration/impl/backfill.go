package impl

import (
	"context"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	coredatabase "grovehub/libs/core/database"
	authorityconstants "grovehub/services/identity/internal/app/repository/authority/constants"
	accountroleimpl "grovehub/services/identity/internal/app/repository/accountrole/impl"
	"grovehub/services/identity/internal/app/repository/rbacmigration/constants"
	"grovehub/services/identity/internal/app/services/rbacmigration"
	"grovehub/services/identity/internal/app/services/rbacmigration/operations/plan"
)

type Backfill struct {
	db *gorm.DB
}

var _ rbacmigration.Backfill = (*Backfill)(nil)

func NewBackfill(db *gorm.DB) rbacmigration.Backfill {
	return &Backfill{db: db}
}

func (writer *Backfill) Apply(ctx context.Context, reviewed rbacmigration.Plan, operatorID string) (rbacmigration.BackfillResult, error) {
	if strings.TrimSpace(reviewed.TenantID) == "" || strings.TrimSpace(reviewed.TenantID) != reviewed.TenantID {
		return rbacmigration.BackfillResult{}, rbacmigration.ErrScopeRequired
	}
	operator, err := uuid.Parse(operatorID)
	if err != nil || operator == uuid.Nil || writer.db == nil {
		return rbacmigration.BackfillResult{}, rbacmigration.ErrInvalidInventory
	}
	var result rbacmigration.BackfillResult
	// The caller must use a fresh operator context, not a nested authority transaction.
	// The runtime must be stopped before an operator calls this primitive.
	err = coredatabase.RunTransaction(ctx, writer.db, func(tx context.Context) error {
		db := coredatabase.TransactionDB(tx, writer.db)
		// Offline only; caller must not already hold tenant authority locks.
		if err := db.Exec(authorityconstants.LockGlobalExclusive).Error; err != nil {
			return err
		}
		if err := db.Exec(constants.LockTables).Error; err != nil {
			return err
		}
		inventory := rbacmigration.Inventory{TenantID: reviewed.TenantID}
		if err := loadInventory(db, reviewed.TenantID, &inventory); err != nil {
			return err
		}
		current, err := plan.Build(inventory)
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(current, reviewed) {
			return rbacmigration.ErrStalePlan
		}
		if err := insertRoles(db, current, operator.String()); err != nil {
			return err
		}
		assignments, err := insertAssignments(db, current, operator.String())
		if err != nil {
			return err
		}
		if err := writer.verify(tx, current, inventory, assignments); err != nil {
			return err
		}
		result = rbacmigration.BackfillResult{RolesCreated: len(current.Roles), AssignmentsCreated: assignments, UnassignedAccounts: len(current.UnassignedAccountIDs)}
		return nil
	})
	if err != nil {
		return rbacmigration.BackfillResult{}, err
	}
	return result, nil
}

func insertRoles(db *gorm.DB, current rbacmigration.Plan, operatorID string) error {
	for _, item := range current.Roles {
		name := strings.ToLower(item.Name)
		if err := db.Exec(constants.InsertRole, item.RoleID, current.TenantID, item.Name, name, name, operatorID, operatorID).Error; err != nil {
			return err
		}
		if err := db.Exec(constants.InsertPersona, current.TenantID, item.RoleID, item.Persona).Error; err != nil {
			return err
		}
	}
	return nil
}

func insertAssignments(db *gorm.DB, current rbacmigration.Plan, operatorID string) (int, error) {
	count := 0
	for _, item := range current.Assignments {
		if err := db.Exec(constants.InsertAssignment, current.TenantID, item.AccountID, item.PrimaryRoleID, "PRIMARY", operatorID).Error; err != nil {
			return 0, err
		}
		count++
		for _, id := range item.SecondaryRoleIDs {
			if err := db.Exec(constants.InsertAssignment, current.TenantID, item.AccountID, id, "SECONDARY", operatorID).Error; err != nil {
				return 0, err
			}
			count++
		}
	}
	return count, nil
}

func (writer *Backfill) verify(ctx context.Context, current rbacmigration.Plan, before rbacmigration.Inventory, expected int) error {
	db := coredatabase.TransactionDB(ctx, writer.db)
	after := rbacmigration.Inventory{TenantID: current.TenantID}
	if err := loadInventory(db, current.TenantID, &after); err != nil {
		return err
	}
	if !reflect.DeepEqual(before.Accounts, after.Accounts) || after.NormalizedAssignments != int64(expected) || len(after.Roles) != len(before.Roles)+len(current.Roles) {
		return rbacmigration.ErrParity
	}
	for _, item := range current.Roles {
		var count int64
		if err := db.Raw(constants.VerifyRole, current.TenantID, item.RoleID, item.Name, strings.ToLower(item.Name), item.Persona).Scan(&count).Error; err != nil {
			return err
		}
		if count != 1 {
			return rbacmigration.ErrParity
		}
	}
	reader := accountroleimpl.New(writer.db)
	for _, item := range current.Assignments {
		actual, err := reader.Find(ctx, current.TenantID, item.AccountID)
		if err != nil {
			return err
		}
		if actual.Version != item.Version || actual.PrimaryRoleID != item.PrimaryRoleID || !reflect.DeepEqual(actual.SecondaryRoleIDs, item.SecondaryRoleIDs) {
			return rbacmigration.ErrParity
		}
	}
	return nil
}

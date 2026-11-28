package impl

import (
	"context"
	"reflect"

	"github.com/google/uuid"
	"gorm.io/gorm"

	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	authorityconstants "grovehub/services/identity/internal/app/repository/authority/constants"
	accountroleimpl "grovehub/services/identity/internal/app/repository/accountrole/impl"
	authorityimpl "grovehub/services/identity/internal/app/repository/authority/impl"
	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/repository/rbacmigration/constants"
	roleimpl "grovehub/services/identity/internal/app/repository/role/impl"
	rolehelpers "grovehub/services/identity/internal/app/repository/role/helpers"
	"grovehub/services/identity/internal/app/services/rbacmigration"
	manifestvalidator "grovehub/services/identity/internal/app/services/rbacmigration/operations/manifest"
	"grovehub/services/identity/internal/app/services/rbacmigration/operations/plan"
)

type Bootstrap struct {
	db *gorm.DB
}

var _ rbacmigration.Bootstrap = (*Bootstrap)(nil)

func NewBootstrap(db *gorm.DB) rbacmigration.Bootstrap {
	return &Bootstrap{db: db}
}

func (writer *Bootstrap) Apply(ctx context.Context, request rbacmigration.Manifest, operatorID string) (rbacmigration.BootstrapResult, error) {
	operator, err := uuid.Parse(operatorID)
	if err != nil || operator == uuid.Nil || writer.db == nil {
		return rbacmigration.BootstrapResult{}, rbacmigration.ErrInvalidManifest
	}
	if err := manifestvalidator.Validate(request); err != nil {
		return rbacmigration.BootstrapResult{}, err
	}
	var result rbacmigration.BootstrapResult
	err = transaction.Run(ctx, writer.db, func(tx context.Context) error {
		db := transaction.Database(tx, writer.db)
		// Offline only; caller must not already hold tenant authority locks.
		if err := db.Exec(authorityconstants.LockGlobalExclusive).Error; err != nil {
			return err
		}
		if err := db.Exec(constants.LockTables).Error; err != nil {
			return err
		}
		if err := db.Exec(constants.LockRegistry).Error; err != nil {
			return err
		}
		if err := writer.verifyBackfill(tx, request.Plan); err != nil {
			return err
		}
		admin, err := accountroleimpl.New(writer.db).Find(tx, request.Plan.TenantID, request.AdministratorID)
		if err != nil {
			return err
		}
		if admin.Status != entity.Active {
			return rbacmigration.ErrAdministrator
		}
		reader := authorityimpl.NewAdministratorReader(writer.db)
		count, err := reader.Count(tx, request.Plan.TenantID)
		if err != nil {
			return err
		}
		if count != 0 {
			return rbacmigration.ErrConflict
		}
		grants, err := writer.activate(tx, request, operator.String())
		if err != nil {
			return err
		}
		count, err = reader.Count(tx, request.Plan.TenantID)
		if err != nil {
			return err
		}
		if count != 1 {
			return rbacmigration.ErrAdministrator
		}
		result = rbacmigration.BootstrapResult{AdministratorID: request.AdministratorID, RolesActivated: len(request.Grants), PermissionsGranted: grants}
		return nil
	})
	if err != nil {
		return rbacmigration.BootstrapResult{}, err
	}
	return result, nil
}

func (writer *Bootstrap) verifyBackfill(ctx context.Context, reviewed rbacmigration.Plan) error {
	db := transaction.Database(ctx, writer.db)
	before := rbacmigration.Inventory{TenantID: reviewed.TenantID}
	if err := loadInventory(db, reviewed.TenantID, &before); err != nil {
		return err
	}
	// Reconstruct the pre-backfill inventory, excluding only the planned roles.
	ids := make(map[string]bool)
	for _, item := range reviewed.Roles {
		ids[item.RoleID] = true
	}
	var existing []rbacmigration.ExistingRole
	for _, item := range before.Roles {
		if !ids[item.RoleID] {
			existing = append(existing, item)
		}
	}
	before.Roles = existing
	before.NormalizedAssignments = 0
	rebuilt, err := plan.Build(before)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(rebuilt, reviewed) {
		return rbacmigration.ErrStalePlan
	}
	expected := 0
	for _, item := range reviewed.Assignments {
		expected += 1 + len(item.SecondaryRoleIDs)
	}
	return (&Backfill{db: writer.db}).verify(ctx, reviewed, before, expected)
}

func (writer *Bootstrap) activate(ctx context.Context, request rbacmigration.Manifest, operatorID string) (int, error) {
	db := transaction.Database(ctx, writer.db)
	repository := roleimpl.New(writer.db)
	total := 0
	for _, grant := range request.Grants {
		detail, err := repository.FindByID(ctx, request.Plan.TenantID, grant.RoleID)
		if err != nil {
			return 0, err
		}
		if err := rolehelpers.CheckActivation(db, detail, grant.Codes); err != nil {
			return 0, err
		}
		ids, err := rolehelpers.ResolvePermissionIDs(db, grant.Codes)
		if err != nil {
			return 0, err
		}
		for _, id := range ids {
			if err := db.Exec(constants.InsertGrant, request.Plan.TenantID, grant.RoleID, id, operatorID).Error; err != nil {
				return 0, err
			}
			total++
		}
		result := db.Exec(constants.ActivateRole, operatorID, request.Plan.TenantID, grant.RoleID)
		if result.Error != nil {
			return 0, result.Error
		}
		if result.RowsAffected != 1 {
			return 0, rbacmigration.ErrStalePlan
		}
	}
	return total, nil
}

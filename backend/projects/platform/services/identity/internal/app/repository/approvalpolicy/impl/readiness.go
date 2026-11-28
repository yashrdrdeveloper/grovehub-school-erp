package impl

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"

	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	kind "grovehub/services/identity/internal/app/constants/approval"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/repository/approval"
	"grovehub/services/identity/internal/app/repository/approvalpolicy"
	"grovehub/services/identity/internal/app/repository/approvalpolicy/constants"
	"grovehub/services/identity/internal/app/repository/authority"
	authorityimpl "grovehub/services/identity/internal/app/repository/authority/impl"
	"grovehub/services/identity/internal/app/repository/rbacmigration"
	migrationimpl "grovehub/services/identity/internal/app/repository/rbacmigration/impl"
	approvalhelpers "grovehub/services/identity/internal/app/services/approval/helpers"
)

type Readiness struct {
	db *gorm.DB
	locker authority.Locker
	migration rbacmigration.Readiness
	administrators authority.AdministratorReader
}

var _ rbacmigration.Readiness = (*Readiness)(nil)

// NewReadiness composes migration and policy readiness. It is not a cutover
// authorization or a substitute for per-effect permission and version checks.
func NewReadiness(db *gorm.DB) rbacmigration.Readiness {
	return &Readiness{db: db, locker: authorityimpl.New(db), migration: migrationimpl.NewReadiness(db), administrators: authorityimpl.NewAdministratorReader(db)}
}

func (reader *Readiness) Require(ctx context.Context, tenantID string) error {
	if reader.db == nil {
		return transaction.ErrMissingDatabase
	}
	return reader.locker.Run(ctx, tenantID, func(tx context.Context) error {
		if err := reader.migration.Require(tx, tenantID); err != nil {
			return err
		}
		count, err := reader.administrators.Count(tx, tenantID)
		if err != nil {
			return err
		}
		if count < 1 {
			return rbacmigration.ErrNotReady
		}
		return reader.requirePolicies(tx, tenantID)
	})
}

func (reader *Readiness) requirePolicies(ctx context.Context, tenantID string) error {
	db := transaction.Database(ctx, reader.db)
	var payload string
	result := db.Raw(constants.LatestReceipt, tenantID).Scan(&payload)
	if result.Error != nil {
		return result.Error
	}
	var receipt approvalpolicy.Receipt
	if result.RowsAffected != 1 || json.Unmarshal([]byte(payload), &receipt) != nil || len(receipt.Policies) != 4 {
		return rbacmigration.ErrNotReady
	}
	expected := map[kind.RequestType]bool{kind.AccountLifecycleChange: true, kind.RoleLifecycleChange: true, kind.RolePermissionChange: true, kind.RoleAssignmentChange: true}
	for _, revision := range receipt.Policies {
		if !expected[revision.Type] || revision.Version < 1 {
			return rbacmigration.ErrNotReady
		}
		delete(expected, revision.Type)
		if err := requireRevision(db, tenantID, revision); err != nil {
			return err
		}
	}
	return nil
}

func requireRevision(db *gorm.DB, tenantID string, revision approvalpolicy.Activated) error {
	var active approvalpolicy.Activated
	result := db.Raw(constants.ActiveRevision, tenantID, revision.Type).Scan(&active)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 || active != revision {
		return rbacmigration.ErrNotReady
	}
	policy := approval.Policy{PolicyID: active.PolicyID}
	if err := db.Raw(constants.ActiveSteps, tenantID, active.PolicyID).Scan(&policy.Steps).Error; err != nil {
		return err
	}
	if err := approvalhelpers.ValidatePolicy(policy); err != nil {
		return rbacmigration.ErrNotReady
	}
	for _, step := range policy.Steps {
		definition, known := permission.Find(step.Permission)
		if !known {
			return rbacmigration.ErrNotReady
		}
		var count int64
		if err := db.Raw(constants.Eligible, tenantID, step.Permission, definition.Feature, definition.Feature, step.RoleID, step.RoleID).Scan(&count).Error; err != nil {
			return err
		}
		if count < 2 {
			return rbacmigration.ErrNotReady
		}
	}
	return nil
}

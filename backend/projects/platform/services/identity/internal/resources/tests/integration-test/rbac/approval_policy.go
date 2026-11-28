package rbac_test

import (
	"context"
	"strings"

	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	coredatabase "grovehub/libs/core/database"
	kind "grovehub/services/identity/internal/app/constants/approval"
	"grovehub/services/identity/internal/app/constants/feature"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/repository/approvalpolicy"
	policyimpl "grovehub/services/identity/internal/app/repository/approvalpolicy/impl"
	"grovehub/services/identity/internal/app/repository/authority"
	authorityimpl "grovehub/services/identity/internal/app/repository/authority/impl"
	"grovehub/services/identity/internal/app/repository/rbacmigration"
	migrationimpl "grovehub/services/identity/internal/app/repository/rbacmigration/impl"
)

func verifyPolicyProvisioning(ctx context.Context, db *gorm.DB, tenant string, author string, reviewerRole string, requestID string, oldPolicyID string) {
	readiness := policyimpl.NewReadiness(db)
	Expect(readiness.Require(ctx, tenant)).To(MatchError(rbacmigration.ErrNotReady))
	provisioner := policyimpl.New(db, migrationimpl.NewReadiness(db), authorityimpl.NewAdministratorReader(db))
	steps := []approvalpolicy.Step{{Permission: permission.AccountRoleAssign}}
	manifest := approvalpolicy.Manifest{TenantID: tenant, AuthorID: author, Policies: []approvalpolicy.Revision{
		{Type: kind.AccountLifecycleChange, ExpectedVersion: 0, Steps: steps},
		{Type: kind.RoleLifecycleChange, ExpectedVersion: 1, Steps: steps},
		{Type: kind.RolePermissionChange, ExpectedVersion: 1, Steps: steps},
		{Type: kind.RoleAssignmentChange, ExpectedVersion: 2, Steps: steps},
	}}
	// Only one account carries this permission in the requester role.
	manifest.Policies[3].Steps = []approvalpolicy.Step{{Permission: permission.RolePermissionManage}}
	_, err := provisioner.Apply(ctx, manifest, author, strings.Repeat("e", 64))
	Expect(err).To(MatchError(approvalpolicy.ErrReviewers))
	var count int64
	Expect(db.Raw("SELECT count(*) FROM identity.approval_policies WHERE tenant_id = ? AND request_type = 'ACCOUNT_LIFECYCLE_CHANGE'", tenant).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(BeZero())
	Expect(db.Raw("SELECT count(*) FROM identity.approval_policy_receipts WHERE tenant_id = ?", tenant).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(BeZero())
	manifest.Policies[3].Steps = steps
	receipt, err := provisioner.Apply(ctx, manifest, author, strings.Repeat("e", 64))
	Expect(err).NotTo(HaveOccurred())
	Expect(receipt.Policies).To(HaveLen(4))
	Expect(readiness.Require(ctx, tenant)).To(Succeed())
	Expect(db.Raw("SELECT count(*) FROM identity.approval_policies WHERE tenant_id = ? AND active", tenant).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(Equal(int64(4)))
	Expect(db.Raw("SELECT count(*) FROM identity.approval_policies WHERE tenant_id = ? AND policy_id = ? AND NOT active", tenant, oldPolicyID).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(Equal(int64(1)))
	Expect(db.Raw("SELECT count(*) FROM identity.approval_requests r JOIN identity.approval_request_steps s ON s.tenant_id = r.tenant_id AND s.request_id = r.request_id WHERE r.tenant_id = ? AND r.request_id = ? AND r.policy_id = ? AND s.required_role_id = ?", tenant, requestID, oldPolicyID, reviewerRole).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(Equal(int64(2)))
	_, err = provisioner.Apply(ctx, manifest, author, strings.Repeat("f", 64))
	Expect(err).To(MatchError(approvalpolicy.ErrConflict))
	Expect(db.Exec("UPDATE identity.approval_policy_receipts SET result = '{}'::jsonb WHERE receipt_id = ?", receipt.ReceiptID).Error).To(HaveOccurred())
	Expect(db.Exec("DELETE FROM identity.approval_policy_receipts WHERE receipt_id = ?", receipt.ReceiptID).Error).To(HaveOccurred())
	// The latest receipt must match the live active revisions.
	Expect(db.Exec("UPDATE identity.approval_policies SET active = false WHERE tenant_id = ? AND policy_id = ?", tenant, receipt.Policies[0].PolicyID).Error).NotTo(HaveOccurred())
	Expect(readiness.Require(ctx, tenant)).To(MatchError(rbacmigration.ErrNotReady))
	Expect(db.Exec("UPDATE identity.approval_policies SET active = true WHERE tenant_id = ? AND policy_id = ?", tenant, receipt.Policies[0].PolicyID).Error).NotTo(HaveOccurred())
	// Eligibility is checked again even when the provisioning receipt exists.
	Expect(db.Exec("UPDATE identity.roles SET status = 'INACTIVE' WHERE tenant_id = ? AND role_id = ?", tenant, reviewerRole).Error).NotTo(HaveOccurred())
	Expect(readiness.Require(ctx, tenant)).To(MatchError(rbacmigration.ErrNotReady))
	Expect(db.Exec("UPDATE identity.roles SET status = 'ACTIVE' WHERE tenant_id = ? AND role_id = ?", tenant, reviewerRole).Error).NotTo(HaveOccurred())
	Expect(readiness.Require(ctx, tenant)).To(Succeed())

	global := authorityimpl.NewGlobalGuard(db)
	administrators := authorityimpl.NewAdministratorReader(db)
	// Removing a global feature, a recovery permission or one tenant entitlement
	// must roll back if any protected school's administrator would disappear.
	for _, statement := range []string{
		"UPDATE identity.features SET active = false WHERE code = ?",
		"UPDATE identity.permissions SET active = false WHERE code = ?",
		"UPDATE identity.tenant_features SET enabled = false WHERE tenant_id = ? AND feature_id = (SELECT feature_id FROM identity.features WHERE code = 'IDENTITY_ROLE_MANAGEMENT')",
	} {
		var argument any = feature.RoleManagement
		if strings.Contains(statement, "identity.permissions") {
			argument = permission.RolePermissionManage
		} else if strings.Contains(statement, "identity.tenant_features") {
			argument = tenant
		}
		Expect(global.Run(ctx, func(tx context.Context) error {
			return coredatabase.TransactionDB(tx, db).Exec(statement, argument).Error
		})).To(MatchError(authority.ErrLastAdministrator))
		count, err := administrators.Count(ctx, tenant)
		Expect(err).NotTo(HaveOccurred())
		Expect(count).To(BeNumerically(">=", 1))
		Expect(readiness.Require(ctx, tenant)).To(Succeed())
	}
	Expect(authorityimpl.New(db).Run(ctx, tenant, func(tx context.Context) error {
		return global.Run(tx, func(context.Context) error {
			return nil
		})
	})).To(MatchError(authority.ErrLockOrder))
	Expect(global.Run(ctx, func(tx context.Context) error {
		return nil
	})).To(Succeed())
}

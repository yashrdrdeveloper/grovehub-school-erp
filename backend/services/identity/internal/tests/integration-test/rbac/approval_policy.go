package rbac_test

import (
	"context"
	"strings"

	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	kind "grovehub/services/identity/internal/app/constants/approval"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/repository/approvalpolicy"
	policyimpl "grovehub/services/identity/internal/app/repository/approvalpolicy/impl"
	authorityimpl "grovehub/services/identity/internal/app/repository/authority/impl"
	migrationimpl "grovehub/services/identity/internal/app/repository/rbacmigration/impl"
)

func verifyPolicyProvisioning(ctx context.Context, db *gorm.DB, tenant string, author string, reviewerRole string, requestID string, oldPolicyID string) {
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
}

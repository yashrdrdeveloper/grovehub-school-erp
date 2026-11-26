package rbac_test

import (
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"grovehub/services/identity/internal/app/constants/permission"
)

func verifyApprovalSchema(db *gorm.DB) {
	const tenant = "approval-school"
	const requester = "10000000-0000-4000-8000-000000000001"
	const reviewer = "10000000-0000-4000-8000-000000000002"
	const foreign = "10000000-0000-4000-8000-000000000003"
	const policy = "10000000-0000-4000-8000-000000000004"
	const request = "10000000-0000-4000-8000-000000000005"
	Expect(db.Exec("INSERT INTO identity.accounts(account_id, person_gid, tenant_id, email, status) VALUES (?, ?, ?, 'requester@approval.test', 'ACTIVE'), (?, ?, ?, 'reviewer@approval.test', 'ACTIVE'), (?, ?, 'foreign-approval-school', 'foreign@approval.test', 'ACTIVE')", requester, requester, tenant, reviewer, reviewer, tenant, foreign, foreign).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.approval_policies(policy_id, tenant_id, request_type, version, created_by) VALUES (?, ?, 'ROLE_ASSIGNMENT_CHANGE', 1, ?)", policy, tenant, requester).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.approval_policy_steps(tenant_id, policy_id, step_number, required_permission_code) VALUES (?, ?, 1, ?)", tenant, policy, permission.AccountRoleAssign).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.approval_requests(request_id, tenant_id, policy_id, requester_id, request_type, payload, payload_sha256, reason, idempotency_key, step_count, expires_at) VALUES (?, ?, ?, ?, 'ROLE_ASSIGNMENT_CHANGE', '{}'::jsonb, repeat('a',64), 'Review change', 'key', 1, now()+interval '1 hour')", request, tenant, policy, requester).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.approval_request_steps(tenant_id, request_id, step_number, required_permission_code) VALUES (?, ?, 1, ?)", tenant, request, permission.AccountRoleAssign).Error).NotTo(HaveOccurred())
	insert := "INSERT INTO identity.approval_decisions(tenant_id, request_id, requester_id, step_number, reviewer_id, decision, comment) VALUES (?, ?, ?, 1, ?, ?, ?)"
	Expect(db.Exec(insert, tenant, request, requester, requester, "APPROVE", "").Error).To(HaveOccurred())
	Expect(db.Exec(insert, tenant, request, requester, foreign, "APPROVE", "").Error).To(HaveOccurred())
	Expect(db.Exec(insert, tenant, request, requester, reviewer, "REJECT", "").Error).To(HaveOccurred())
	Expect(db.Exec(insert, tenant, request, requester, reviewer, "APPROVE", "").Error).NotTo(HaveOccurred())
	Expect(db.Exec(insert, tenant, request, requester, reviewer, "APPROVE", "").Error).To(HaveOccurred())
	Expect(db.Exec("UPDATE identity.approval_decisions SET comment = 'changed' WHERE tenant_id = ? AND request_id = ?", tenant, request).Error).To(HaveOccurred())
	Expect(db.Exec("DELETE FROM identity.approval_decisions WHERE tenant_id = ? AND request_id = ?", tenant, request).Error).To(HaveOccurred())
	var count int64
	Expect(db.Raw("SELECT count(*) FROM identity.approval_decisions WHERE tenant_id = ? AND request_id = ?", tenant, request).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(Equal(int64(1)))
}

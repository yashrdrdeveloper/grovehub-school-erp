package rbac_test

import (
	"context"
	"time"

	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"grovehub/libs/schoolgateway"
	kind "grovehub/services/identity/internal/app/constants/approval"
	"grovehub/services/identity/internal/app/constants/feature"
	"grovehub/services/identity/internal/app/constants/permission"
	dto "grovehub/services/identity/internal/app/dto/approval"
	"grovehub/services/identity/internal/app/entity"
	accountimpl "grovehub/services/identity/internal/app/repository/account/impl"
	repo "grovehub/services/identity/internal/app/repository/approval"
	repository "grovehub/services/identity/internal/app/repository/approval/impl"
	authorityimpl "grovehub/services/identity/internal/app/repository/authority/impl"
	authrepo "grovehub/services/identity/internal/app/repository/authorization/impl"
	registryrepo "grovehub/services/identity/internal/app/repository/registry/impl"
	roleimpl "grovehub/services/identity/internal/app/repository/role/impl"
	approvalimpl "grovehub/services/identity/internal/app/services/approval/impl"
	"grovehub/services/identity/internal/app/services/approval/operations/state"
	"grovehub/services/identity/internal/app/services/authorization"
	authimpl "grovehub/services/identity/internal/app/services/authorization/impl"
	"grovehub/services/identity/internal/app/services/registry"
	registryimpl "grovehub/services/identity/internal/app/services/registry/impl"
)

func verifyApprovalWorkflow(ctx context.Context, db *gorm.DB) {
	const tenant = "approval-school"
	const requester = "10000000-0000-4000-8000-000000000001"
	const reviewer = "10000000-0000-4000-8000-000000000002"
	const foreign = "10000000-0000-4000-8000-000000000003"
	const policyID = "10000000-0000-4000-8000-000000000006"
	const managementRole = "10000000-0000-4000-8000-000000000007"
	const targetRole = "10000000-0000-4000-8000-000000000008"
	const requiredRole = "10000000-0000-4000-8000-000000000009"

	catalog := registryimpl.New(registryrepo.New(db))
	_, err := catalog.Register(ctx, registry.Request{Features: []feature.Code{feature.RoleManagement}, Permissions: []permission.Code{permission.RoleDeactivate, permission.AccountRoleAssign}})
	Expect(err).NotTo(HaveOccurred())
	Expect(db.Exec("UPDATE identity.features SET active = true WHERE code = ?", feature.RoleManagement).Error).NotTo(HaveOccurred())
	Expect(db.Exec("UPDATE identity.permissions SET active = true WHERE code IN ?", []permission.Code{permission.RoleDeactivate, permission.AccountRoleAssign}).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.tenant_features(tenant_id, feature_id, enabled) SELECT ?, feature_id, true FROM identity.features WHERE code = ?", tenant, feature.RoleManagement).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.roles(role_id, tenant_id, name, normalized_name, status) VALUES (?, ?, 'Requester authority', 'requester authority', 'ACTIVE'), (?, ?, 'Target', 'target', 'ACTIVE'), (?, ?, 'Reviewer authority', 'reviewer authority', 'ACTIVE')", managementRole, tenant, targetRole, tenant, requiredRole, tenant).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.role_permissions(tenant_id, role_id, permission_id) SELECT ?, ?, permission_id FROM identity.permissions WHERE code = ?", tenant, managementRole, permission.RoleDeactivate).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.role_permissions(tenant_id, role_id, permission_id) SELECT ?, ?, permission_id FROM identity.permissions WHERE code = ?", tenant, requiredRole, permission.AccountRoleAssign).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.account_roles(tenant_id, account_id, role_id, assignment_type) VALUES (?, ?, ?, 'PRIMARY'), (?, ?, ?, 'PRIMARY')", tenant, requester, managementRole, tenant, reviewer, requiredRole).Error).NotTo(HaveOccurred())
	snapshots := authrepo.New(db)
	service := approvalimpl.New(repository.New(db), authorityimpl.New(db), authimpl.New(snapshots), snapshots, accountimpl.NewNormalizedReader(db), roleimpl.New(db))
	requesterCtx := schoolgateway.WithPrincipal(ctx, schoolgateway.Principal{TenantID: tenant, AccountID: requester, SessionID: "requester-session"})
	reviewerCtx := schoolgateway.WithPrincipal(ctx, schoolgateway.Principal{TenantID: tenant, AccountID: reviewer, SessionID: "reviewer-session"})
	foreignCtx := schoolgateway.WithPrincipal(ctx, schoolgateway.Principal{TenantID: "foreign-approval-school", AccountID: foreign, SessionID: "foreign-session"})
	intent := dto.CreateLifecycle{Type: kind.RoleLifecycleChange, Intent: dto.LifecycleIntent{TargetID: targetRole, Version: 1, Action: "deactivate"}, Reason: "Review role retirement", IdempotencyKey: "workflow"}
	_, err = service.CreateLifecycle(requesterCtx, intent)
	Expect(err).To(MatchError(repo.ErrPolicy))
	Expect(db.Exec("INSERT INTO identity.approval_policies(policy_id, tenant_id, request_type, version, active, created_by) VALUES (?, ?, 'ROLE_LIFECYCLE_CHANGE', 1, true, ?)", policyID, tenant, requester).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.approval_policy_steps(tenant_id, policy_id, step_number, required_permission_code, required_role_id) VALUES (?, ?, 1, ?, ?), (?, ?, 2, ?, ?)", tenant, policyID, permission.AccountRoleAssign, requiredRole, tenant, policyID, permission.AccountRoleAssign, requiredRole).Error).NotTo(HaveOccurred())
	created, err := service.CreateLifecycle(requesterCtx, intent)
	Expect(err).NotTo(HaveOccurred())
	Expect(created.Status).To(Equal(kind.Pending))
	Expect(created.StepCount).To(Equal(2))
	replayed, err := service.CreateLifecycle(requesterCtx, intent)
	Expect(err).NotTo(HaveOccurred())
	Expect(replayed.RequestID).To(Equal(created.RequestID))
	Expect(replayed.PayloadSHA256).To(Equal(created.PayloadSHA256))
	conflicting := intent
	conflicting.Reason = "Different reason"
	_, err = service.CreateLifecycle(requesterCtx, conflicting)
	Expect(err).To(MatchError(repo.ErrConflict))
	// Later policy changes cannot rewrite the snapshotted eligibility criteria.
	Expect(db.Exec("UPDATE identity.approval_policy_steps SET required_permission_code = ? WHERE tenant_id = ? AND policy_id = ?", permission.RoleDeactivate, tenant, policyID).Error).NotTo(HaveOccurred())
	review := dto.Review{RequestID: created.RequestID, Version: 1, Decision: kind.Approve}
	_, err = service.Decide(requesterCtx, review)
	Expect(err).To(MatchError(state.ErrForbidden))
	_, err = service.Decide(foreignCtx, review)
	Expect(err).To(MatchError(repo.ErrNotFound))
	Expect(db.Exec("UPDATE identity.roles SET status = 'INACTIVE' WHERE tenant_id = ? AND role_id = ?", tenant, requiredRole).Error).NotTo(HaveOccurred())
	_, err = service.Decide(reviewerCtx, review)
	Expect(err).To(MatchError(authorization.ErrForbidden))
	Expect(db.Exec("UPDATE identity.roles SET status = 'ACTIVE' WHERE tenant_id = ? AND role_id = ?", tenant, requiredRole).Error).NotTo(HaveOccurred())
	outcomes := make(chan error, 2)
	for range 2 {
		go func() {
			_, err := service.Decide(reviewerCtx, review)
			outcomes <- err
		}()
	}
	successes := 0
	for range 2 {
		err := <-outcomes
		if err == nil {
			successes++
		} else {
			Expect(err).To(MatchError(state.ErrVersion))
		}
	}
	Expect(successes).To(Equal(1))
	review.Version = 2
	approved, err := service.Decide(reviewerCtx, review)
	Expect(err).NotTo(HaveOccurred())
	Expect(approved.Status).To(Equal(kind.Approved))
	Expect(approved.Version).To(Equal(int64(3)))
	var count int64
	Expect(db.Raw("SELECT count(*) FROM identity.approval_decisions WHERE tenant_id = ? AND request_id = ?", tenant, created.RequestID).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(Equal(int64(2)))
	unchanged, err := roleimpl.New(db).FindByID(ctx, tenant, targetRole)
	Expect(err).NotTo(HaveOccurred())
	Expect(unchanged.Role.Status).To(Equal(entity.RoleActive))
	Expect(unchanged.Role.Version).To(Equal(int64(1)))
	intent.IdempotencyKey = "cancel-workflow"
	cancelRequest, err := service.CreateLifecycle(requesterCtx, intent)
	Expect(err).NotTo(HaveOccurred())
	_, err = service.Cancel(reviewerCtx, cancelRequest.RequestID, 1)
	Expect(err).To(MatchError(state.ErrForbidden))
	cancelled, err := service.Cancel(requesterCtx, cancelRequest.RequestID, 1)
	Expect(err).NotTo(HaveOccurred())
	Expect(cancelled.Status).To(Equal(kind.Cancelled))
	intent.IdempotencyKey = "expired-workflow"
	expired, err := service.CreateLifecycle(requesterCtx, intent)
	Expect(err).NotTo(HaveOccurred())
	Expect(db.Exec("UPDATE identity.approval_requests SET created_at = ?, expires_at = ? WHERE request_id = ?", time.Now().UTC().Add(-2*time.Hour), time.Now().UTC().Add(-time.Hour), expired.RequestID).Error).NotTo(HaveOccurred())
	_, err = service.Decide(reviewerCtx, dto.Review{RequestID: expired.RequestID, Version: 1, Decision: kind.Approve})
	Expect(err).To(MatchError(state.ErrExpired))
}

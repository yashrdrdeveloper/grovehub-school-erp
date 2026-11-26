package rbac_test

import (
	"context"
	"encoding/json"
	headercontext "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/headercontext"
	"net/http"
	"net/http/httptest"
	"strings"
	approvalcontroller "grovehub/services/identity/internal/app/controller/approval/protected/impl"
	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
	"errors"
	"time"

	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"grovehub/libs/schoolgateway"
	infrastructure "grovehub/services/infrastructure"
	messaging "grovehub/services/infrastructure/client/messaging"
	auditimpl "grovehub/services/identity/internal/app/repository/approvalaudit/impl"
	resourceaudit "grovehub/services/identity/internal/app/repository/resourceaudit/impl"
	resourceproducer "grovehub/services/identity/internal/app/events/producers/resource/impl"
	producerimpl "grovehub/services/identity/internal/app/events/producers/approval/impl"
	coreevents "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events"
	kind "grovehub/services/identity/internal/app/constants/approval"
	"grovehub/services/identity/internal/app/constants/feature"
	"grovehub/services/identity/internal/app/constants/permission"
	dto "grovehub/services/identity/internal/app/dto/approval"
	"grovehub/services/identity/internal/app/entity"
	accountdto "grovehub/services/identity/internal/app/dto/accountrole"
	assignmentrepo "grovehub/services/identity/internal/app/repository/accountrole"
	assignmentimpl "grovehub/services/identity/internal/app/repository/accountrole/impl"
	accountimpl "grovehub/services/identity/internal/app/repository/account/impl"
	repo "grovehub/services/identity/internal/app/repository/approval"
	repository "grovehub/services/identity/internal/app/repository/approval/impl"
	authorityimpl "grovehub/services/identity/internal/app/repository/authority/impl"
	"grovehub/services/identity/internal/app/repository/authority"
	migrationimpl "grovehub/services/identity/internal/app/repository/rbacmigration/impl"
	migrationrepo "grovehub/services/identity/internal/app/repository/rbacmigration"
	rolerepo "grovehub/services/identity/internal/app/repository/role"
	"grovehub/services/identity/internal/app/services/authorization/operations/protect"
	"grovehub/services/identity/internal/app/services/accountrole/operations/replace"
	"grovehub/services/identity/internal/app/services/role/operations/lifecycle"
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
	store := infrastructure.NewEventStore(db)
	recorder := producerimpl.New(auditimpl.New(db), store)
	service := approvalimpl.New(repository.New(db, recorder), authorityimpl.New(db), authimpl.New(snapshots), snapshots, accountimpl.NewNormalizedReader(db), roleimpl.New(db), assignmentimpl.New(db))
	requesterCtx := schoolgateway.WithPrincipal(coreevents.WithCorrelationID(ctx, "approval-audit-test"), schoolgateway.Principal{TenantID: tenant, AccountID: requester, SessionID: "requester-session"})
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
	reads := approvalimpl.NewReader(repository.NewReader(db), repository.New(db, recorder), authorityimpl.New(db), snapshots)
	query := coresearch.Request{PageInfo: coresearch.PageRequest{PageNumber: 1, RowSize: 25}}
	historyQuery := query
	historyQuery.Filters = []coresearch.Filter{{Field: dto.FieldRequestID, Operator: coresearch.OperatorEquals, Values: []string{created.RequestID}}}
	history, err := reads.History(requesterCtx, historyQuery)
	Expect(err).NotTo(HaveOccurred())
	Expect(history.Items).To(HaveLen(1))
	Expect(history.Items[0].RequestID).To(Equal(created.RequestID))
	inbox, err := reads.Inbox(reviewerCtx, query)
	Expect(err).NotTo(HaveOccurred())
	Expect(inbox.Items).To(HaveLen(1))
	Expect(inbox.Items[0].RequestID).To(Equal(created.RequestID))
	pendingCount, err := reads.PendingCount(reviewerCtx)
	Expect(err).NotTo(HaveOccurred())
	Expect(pendingCount).To(Equal(int64(1)))
	pendingCount, err = reads.PendingCount(requesterCtx)
	Expect(err).NotTo(HaveOccurred())
	Expect(pendingCount).To(BeZero())
	foreignInbox, err := reads.Inbox(foreignCtx, query)
	Expect(err).NotTo(HaveOccurred())
	Expect(foreignInbox.Items).To(BeEmpty())
	detail, err := reads.Detail(reviewerCtx, created.RequestID)
	Expect(err).NotTo(HaveOccurred())
	Expect(detail.Reason).To(Equal(intent.Reason))
	controller := approvalcontroller.New(reads)
	httpRequest := httptest.NewRequest(http.MethodGet, "/identity/protected/approvals/"+created.RequestID, nil).WithContext(reviewerCtx)
	httpRequest.SetPathValue("requestId", created.RequestID)
	response := httptest.NewRecorder()
	controller.Detail(response, httpRequest)
	Expect(response.Code).To(Equal(http.StatusOK))
	Expect(response.Body.String()).To(ContainSubstring(created.RequestID))
	Expect(response.Body.String()).NotTo(ContainSubstring("idempotencyKey"))
	response = httptest.NewRecorder()
	controller.Detail(response, httpRequest.WithContext(foreignCtx))
	Expect(response.Code).To(Equal(http.StatusNotFound))
	response = httptest.NewRecorder()
	controller.Inbox(response, httptest.NewRequest(http.MethodPost, "/identity/protected/approvals/inbox/search", strings.NewReader("{")).WithContext(reviewerCtx))
	Expect(response.Code).To(Equal(http.StatusBadRequest))
	_, err = reads.Detail(foreignCtx, created.RequestID)
	Expect(err).To(MatchError(repo.ErrNotFound))
	events, err := reads.Audit(requesterCtx, created.RequestID)
	Expect(err).NotTo(HaveOccurred())
	Expect(events).To(HaveLen(1))
	_, err = reads.Audit(reviewerCtx, created.RequestID)
	Expect(err).To(MatchError(repo.ErrNotFound))
	replayed, err := service.CreateLifecycle(requesterCtx, intent)
	Expect(err).NotTo(HaveOccurred())
	Expect(replayed.RequestID).To(Equal(created.RequestID))
	Expect(replayed.PayloadSHA256).To(Equal(created.PayloadSHA256))
	var auditCount int64
	Expect(db.Raw("SELECT count(*) FROM identity.approval_audit WHERE tenant_id = ? AND request_id = ?", tenant, created.RequestID).Scan(&auditCount).Error).NotTo(HaveOccurred())
	Expect(auditCount).To(Equal(int64(1)))
	Expect(db.Raw("SELECT count(*) FROM identity.approval_audit a JOIN infrastructure.messaging m ON m.event_id = a.event_id::text WHERE a.tenant_id = ? AND a.request_id = ? AND a.event_payload = convert_from(m.payload, 'UTF8')::jsonb AND a.event_payload->>'correlationId' = 'approval-audit-test' AND a.event_payload->>'causationId' = a.request_id::text", tenant, created.RequestID).Scan(&auditCount).Error).NotTo(HaveOccurred())
	Expect(auditCount).To(Equal(int64(1)))
	Expect(db.Exec("UPDATE identity.approval_audit SET operation = 'CANCELLED' WHERE tenant_id = ? AND request_id = ?", tenant, created.RequestID).Error).To(HaveOccurred())
	Expect(db.Exec("DELETE FROM identity.approval_audit WHERE tenant_id = ? AND request_id = ?", tenant, created.RequestID).Error).To(HaveOccurred())
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
	pendingCount, err = reads.PendingCount(reviewerCtx)
	Expect(err).NotTo(HaveOccurred())
	Expect(pendingCount).To(BeZero())
	_, err = reads.Detail(reviewerCtx, created.RequestID)
	Expect(err).To(MatchError(repo.ErrNotFound))
	var count int64
	Expect(db.Raw("SELECT count(*) FROM identity.approval_decisions WHERE tenant_id = ? AND request_id = ?", tenant, created.RequestID).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(Equal(int64(2)))
	unchanged, err := roleimpl.New(db).FindByID(ctx, tenant, targetRole)
	Expect(err).NotTo(HaveOccurred())
	Expect(unchanged.Role.Status).To(Equal(entity.RoleActive))
	Expect(unchanged.Role.Version).To(Equal(int64(1)))
	// Restore the active policy after verifying that existing snapshots did not change.
	Expect(db.Exec("UPDATE identity.approval_policy_steps SET required_permission_code = ? WHERE tenant_id = ? AND policy_id = ?", permission.AccountRoleAssign, tenant, policyID).Error).NotTo(HaveOccurred())
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

	mutationController := approvalcontroller.NewWithMutations(reads, service, nil)
	body, err := json.Marshal(intent)
	Expect(err).NotTo(HaveOccurred())
	submit := func(key string) *httptest.ResponseRecorder {
		request := httptest.NewRequest(http.MethodPost, "/identity/protected/approvals/lifecycle", strings.NewReader(string(body))).WithContext(requesterCtx)
		if key != "" {
			request.Header.Set("Idempotency-Key", key)
		}
		response := httptest.NewRecorder()
		headercontext.Middleware(http.HandlerFunc(mutationController.SubmitLifecycle)).ServeHTTP(response, request)
		return response
	}
	Expect(submit("").Code).To(Equal(http.StatusBadRequest))
	submitted := submit("http-approval-1")
	Expect(submitted.Code).To(Equal(http.StatusCreated))
	persistedHTTP, err := repository.New(db, recorder).FindByKey(ctx, tenant, requester, "http-approval-1")
	Expect(err).NotTo(HaveOccurred())
	Expect(submitted.Body.String()).To(ContainSubstring(persistedHTTP.RequestID))
	Expect(submitted.Body.String()).NotTo(ContainSubstring("idempotencyKey"))
	Expect(submit("http-approval-1").Body.String()).To(ContainSubstring(persistedHTTP.RequestID))
	post := func(handler http.HandlerFunc, actor context.Context, id string, body string) *httptest.ResponseRecorder {
		request := httptest.NewRequest(http.MethodPost, "/identity/protected/approvals/"+id, strings.NewReader(body)).WithContext(actor)
		request.SetPathValue("requestId", id)
		response := httptest.NewRecorder()
		handler(response, request)
		return response
	}
	Expect(post(mutationController.Approve, requesterCtx, persistedHTTP.RequestID, `{"version":1}`).Code).To(Equal(http.StatusForbidden))
	Expect(post(mutationController.Approve, foreignCtx, persistedHTTP.RequestID, `{"version":1}`).Code).To(Equal(http.StatusNotFound))
	Expect(post(mutationController.Approve, reviewerCtx, persistedHTTP.RequestID, `{"version":1}`).Code).To(Equal(http.StatusOK))
	Expect(post(mutationController.Approve, reviewerCtx, persistedHTTP.RequestID, `{"version":1}`).Code).To(Equal(http.StatusConflict))
	Expect(post(mutationController.Reject, reviewerCtx, persistedHTTP.RequestID, `{"version":2,"comment":"Rejected in HTTP test"}`).Code).To(Equal(http.StatusOK))
	Expect(post(mutationController.Execute, requesterCtx, persistedHTTP.RequestID, `{"version":3}`).Code).To(Equal(http.StatusServiceUnavailable))
	Expect(submit("http-approval-2").Code).To(Equal(http.StatusCreated))
	cancelHTTP, err := repository.New(db, recorder).FindByKey(ctx, tenant, requester, "http-approval-2")
	Expect(err).NotTo(HaveOccurred())
	Expect(post(mutationController.Cancel, reviewerCtx, cancelHTTP.RequestID, `{"version":1}`).Code).To(Equal(http.StatusForbidden))
	Expect(post(mutationController.Cancel, requesterCtx, cancelHTTP.RequestID, `{"version":1}`).Code).To(Equal(http.StatusOK))

	const draftID = "10000000-0000-4000-8000-000000000010"
	const permissionPolicy = "10000000-0000-4000-8000-000000000011"
	const assignmentPolicy = "10000000-0000-4000-8000-000000000012"
	Expect(db.Exec("INSERT INTO identity.roles(role_id, tenant_id, name, normalized_name, status) VALUES (?, ?, 'Draft', 'draft', 'INACTIVE')", draftID, tenant).Error).NotTo(HaveOccurred())
	Expect(db.Exec("UPDATE identity.permissions SET active = true WHERE code = ?", permission.RolePermissionManage).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.role_permissions(tenant_id, role_id, permission_id) SELECT ?, ?, permission_id FROM identity.permissions WHERE code IN ?", tenant, managementRole, []permission.Code{permission.RolePermissionManage, permission.AccountRoleAssign}).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.approval_policies(policy_id, tenant_id, request_type, version, active, created_by) VALUES (?, ?, 'ROLE_PERMISSION_CHANGE', 1, true, ?), (?, ?, 'ROLE_ASSIGNMENT_CHANGE', 2, true, ?)", permissionPolicy, tenant, requester, assignmentPolicy, tenant, requester).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.approval_policy_steps(tenant_id, policy_id, step_number, required_permission_code, required_role_id) VALUES (?, ?, 1, ?, ?), (?, ?, 1, ?, ?)", tenant, permissionPolicy, permission.AccountRoleAssign, requiredRole, tenant, assignmentPolicy, permission.AccountRoleAssign, requiredRole).Error).NotTo(HaveOccurred())
	permissionRequest := dto.CreatePermissions{Intent: dto.PermissionIntent{RoleID: draftID, Version: 1, Codes: []permission.Code{permission.RoleDeactivate}}, Reason: "Draft capability", IdempotencyKey: "permission-change"}
	pendingPermission, err := service.CreatePermissions(requesterCtx, permissionRequest)
	Expect(err).NotTo(HaveOccurred())
	Expect(pendingPermission.RequestType).To(Equal(kind.RolePermissionChange))
	approvedPermission, err := service.Decide(reviewerCtx, dto.Review{RequestID: pendingPermission.RequestID, Version: 1, Decision: kind.Approve})
	Expect(err).NotTo(HaveOccurred())
	Expect(approvedPermission.Status).To(Equal(kind.Approved))
	grantSet, err := roleimpl.New(db).Permissions(ctx, tenant, draftID)
	Expect(err).NotTo(HaveOccurred())
	Expect(grantSet.Codes).To(BeEmpty())
	permissionRequest.IdempotencyKey = "tampered-permission"
	tampered, err := service.CreatePermissions(requesterCtx, permissionRequest)
	Expect(err).NotTo(HaveOccurred())
	Expect(db.Exec("UPDATE identity.approval_requests SET payload = jsonb_set(payload, '{version}', '2'::jsonb) WHERE request_id = ?", tampered.RequestID).Error).NotTo(HaveOccurred())
	_, err = service.Decide(reviewerCtx, dto.Review{RequestID: tampered.RequestID, Version: 1, Decision: kind.Approve})
	Expect(err).To(MatchError(repo.ErrConflict))

	assignment := dto.CreateAssignment{Intent: dto.AssignmentIntent{AccountID: reviewer, Version: 1, PrimaryRoleID: managementRole, SecondaryRoleIDs: []string{}, RoleVersions: []accountdto.RoleVersion{{RoleID: managementRole, Version: 1}, {RoleID: requiredRole, Version: 1}}}, Reason: "Change reviewer assignment", IdempotencyKey: "assignment-change"}
	pendingAssignment, err := service.CreateAssignment(requesterCtx, assignment)
	Expect(err).NotTo(HaveOccurred())
	Expect(pendingAssignment.RequestType).To(Equal(kind.RoleAssignmentChange))
	_, err = service.Decide(reviewerCtx, dto.Review{RequestID: pendingAssignment.RequestID, Version: 1, Decision: kind.Approve})
	Expect(err).To(MatchError(state.ErrForbidden))
	assignment.IdempotencyKey = "stale-assignment"
	assignment.Intent.RoleVersions[0].Version++
	_, err = service.CreateAssignment(requesterCtx, assignment)
	Expect(err).To(MatchError(assignmentrepo.ErrVersionConflict))
	assignment.IdempotencyKey = "requester-assignment"
	assignment.Intent.AccountID = requester
	assignment.Intent.PrimaryRoleID = requiredRole
	assignment.Intent.RoleVersions[0].Version--
	pendingAssignment, err = service.CreateAssignment(requesterCtx, assignment)
	Expect(err).NotTo(HaveOccurred())
	approvedAssignment, err := service.Decide(reviewerCtx, dto.Review{RequestID: pendingAssignment.RequestID, Version: 1, Decision: kind.Approve})
	Expect(err).NotTo(HaveOccurred())
	Expect(approvedAssignment.Status).To(Equal(kind.Approved))
	originalAssignment, err := assignmentimpl.New(db).Find(ctx, tenant, requester)
	Expect(err).NotTo(HaveOccurred())
	Expect(originalAssignment.PrimaryRoleID).To(Equal(managementRole))
	Expect(originalAssignment.Version).To(Equal(int64(1)))

	// Approve a second intent against the original draft version before applying either.
	permissionRequest.IdempotencyKey = "stale-after-effect"
	pendingStale, err := service.CreatePermissions(requesterCtx, permissionRequest)
	Expect(err).NotTo(HaveOccurred())
	approvedStale, err := service.Decide(reviewerCtx, dto.Review{RequestID: pendingStale.RequestID, Version: 1, Decision: kind.Approve})
	Expect(err).NotTo(HaveOccurred())
	locker := authorityimpl.New(db)
	guard := protect.New(locker, authorityimpl.NewAdministratorReader(db))
	readiness := migrationimpl.NewReadiness(db)
	effects := approvalimpl.Effects{Audit: resourceproducer.New(resourceaudit.New(db), store), Permissions: roleimpl.NewDraftWriter(db), Assignment: replace.New(guard, assignmentimpl.NewWriter(db)), RoleLifecycle: lifecycle.New(guard, roleimpl.NewLifecycleWriter(db))}
	buildExecutor := func(effects approvalimpl.Effects) *approvalimpl.Executor {
		return approvalimpl.NewExecutor(repository.New(db, recorder), locker, guard, readiness, authimpl.New(snapshots), snapshots, accountimpl.NewNormalizedReader(db), roleimpl.New(db), assignmentimpl.New(db), effects)
	}
	executor := buildExecutor(effects)
	_, err = executor.Execute(requesterCtx, approvedPermission.RequestID, approvedPermission.Version)
	Expect(err).To(MatchError(migrationrepo.ErrNotReady))
	Expect(db.Exec("INSERT INTO identity.rbac_migration_receipts(receipt_id, tenant_id, operator_id, operation, input_sha256, result) VALUES ('10000000-0000-4000-8000-000000000013', ?, ?, 'BACKFILL', repeat('c',64), '{}'::jsonb), ('10000000-0000-4000-8000-000000000014', ?, ?, 'BOOTSTRAP', repeat('d',64), '{}'::jsonb)", tenant, requester, tenant, requester).Error).NotTo(HaveOccurred())
	_, err = executor.Execute(reviewerCtx, approvedPermission.RequestID, approvedPermission.Version)
	Expect(err).To(MatchError(state.ErrForbidden))
	Expect(db.Exec("UPDATE identity.roles SET status = 'INACTIVE' WHERE tenant_id = ? AND role_id = ?", tenant, requiredRole).Error).NotTo(HaveOccurred())
	_, err = executor.Execute(requesterCtx, approvedPermission.RequestID, approvedPermission.Version)
	Expect(err).To(MatchError(authorization.ErrForbidden))
	Expect(db.Exec("UPDATE identity.roles SET status = 'ACTIVE' WHERE tenant_id = ? AND role_id = ?", tenant, requiredRole).Error).NotTo(HaveOccurred())

	fault := effects
	fault.Permissions = failingDraftWriter{DraftPermissionWriter: effects.Permissions}
	_, err = buildExecutor(fault).Execute(requesterCtx, approvedPermission.RequestID, approvedPermission.Version)
	Expect(err).To(MatchError(errInjectedEffect))
	grantSet, err = roleimpl.New(db).Permissions(ctx, tenant, draftID)
	Expect(err).NotTo(HaveOccurred())
	Expect(grantSet.Version).To(Equal(int64(1)))
	Expect(grantSet.Codes).To(BeEmpty())
	var persisted entity.ApprovalRequest
	Expect(db.Raw("SELECT * FROM identity.approval_requests WHERE tenant_id = ? AND request_id = ?", tenant, approvedPermission.RequestID).Scan(&persisted).Error).NotTo(HaveOccurred())
	Expect(persisted.Status).To(Equal(kind.Approved))

	failingRecorder := producerimpl.New(auditimpl.New(db), failingOutbox{})
	failingExecutor := approvalimpl.NewExecutor(repository.New(db, failingRecorder), locker, guard, readiness, authimpl.New(snapshots), snapshots, accountimpl.NewNormalizedReader(db), roleimpl.New(db), assignmentimpl.New(db), effects)
	_, err = failingExecutor.Execute(requesterCtx, approvedPermission.RequestID, approvedPermission.Version)
	Expect(err).To(MatchError(errInjectedOutbox))
	grantSet, err = roleimpl.New(db).Permissions(ctx, tenant, draftID)
	Expect(err).NotTo(HaveOccurred())
	Expect(grantSet.Version).To(Equal(int64(1)))
	Expect(grantSet.Codes).To(BeEmpty())
	Expect(db.Raw("SELECT count(*) FROM identity.approval_audit WHERE tenant_id = ? AND request_id = ? AND operation = 'APPLIED'", tenant, approvedPermission.RequestID).Scan(&auditCount).Error).NotTo(HaveOccurred())
	Expect(auditCount).To(BeZero())
	Expect(db.Raw("SELECT * FROM identity.approval_requests WHERE tenant_id = ? AND request_id = ?", tenant, approvedPermission.RequestID).Scan(&persisted).Error).NotTo(HaveOccurred())
	Expect(persisted.Status).To(Equal(kind.Approved))

	Expect(db.Raw("SELECT count(*) FROM identity.resource_audit WHERE tenant_id = ? AND request_id = ?", tenant, approvedPermission.RequestID).Scan(&auditCount).Error).NotTo(HaveOccurred())
	Expect(auditCount).To(BeZero())
	fault = effects
	fault.Audit = resourceproducer.New(resourceaudit.New(db), failingOutbox{})
	_, err = buildExecutor(fault).Execute(requesterCtx, approvedPermission.RequestID, approvedPermission.Version)
	Expect(err).To(MatchError(errInjectedOutbox))
	grantSet, err = roleimpl.New(db).Permissions(ctx, tenant, draftID)
	Expect(err).NotTo(HaveOccurred())
	Expect(grantSet.Version).To(Equal(int64(1)))
	Expect(grantSet.Codes).To(BeEmpty())
	Expect(db.Raw("SELECT count(*) FROM identity.resource_audit WHERE tenant_id = ? AND request_id = ?", tenant, approvedPermission.RequestID).Scan(&auditCount).Error).NotTo(HaveOccurred())
	Expect(auditCount).To(BeZero())

	applied, err := executor.Execute(requesterCtx, approvedPermission.RequestID, approvedPermission.Version)
	Expect(err).NotTo(HaveOccurred())
	Expect(applied.Status).To(Equal(kind.Applied))
	Expect(applied.Version).To(Equal(approvedPermission.Version+1))
	replay, err := executor.Execute(requesterCtx, approvedPermission.RequestID, approvedPermission.Version)
	Expect(err).NotTo(HaveOccurred())
	Expect(replay.Version).To(Equal(applied.Version))
	Expect(db.Raw("SELECT count(*) FROM identity.approval_audit a JOIN infrastructure.messaging m ON m.event_id = a.event_id::text WHERE a.tenant_id = ? AND a.request_id = ? AND a.operation = 'APPLIED'", tenant, approvedPermission.RequestID).Scan(&auditCount).Error).NotTo(HaveOccurred())
	Expect(auditCount).To(Equal(int64(1)))
	grantSet, err = roleimpl.New(db).Permissions(ctx, tenant, draftID)
	Expect(err).NotTo(HaveOccurred())
	Expect(grantSet.Version).To(Equal(int64(2)))
	Expect(grantSet.Codes).To(Equal([]permission.Code{permission.RoleDeactivate}))
	Expect(db.Raw("SELECT count(*) FROM identity.resource_audit a JOIN infrastructure.messaging m ON m.event_id = a.event_id::text WHERE a.tenant_id = ? AND a.request_id = ? AND a.resource_id = ? AND a.resource_version = 2 AND a.event_payload = convert_from(m.payload, 'UTF8')::jsonb AND a.event_payload->>'eventType' = 'identity.role.changed.v1' AND a.event_payload->>'causationId' = a.request_id::text AND a.event_payload->>'correlationId' = 'approval-audit-test' AND a.event_payload#>'{data,old,role,permissionCodes}' = '[]'::jsonb AND a.event_payload#>'{data,new,role,permissionCodes}' = jsonb_build_array(CAST(? AS text))", tenant, approvedPermission.RequestID, draftID, permission.RoleDeactivate).Scan(&auditCount).Error).NotTo(HaveOccurred())
	Expect(auditCount).To(Equal(int64(1)))
	Expect(db.Exec("UPDATE identity.resource_audit SET operation = 'ROLE_LIFECYCLE_CHANGE' WHERE tenant_id = ? AND request_id = ?", tenant, approvedPermission.RequestID).Error).To(HaveOccurred())
	Expect(db.Exec("DELETE FROM identity.resource_audit WHERE tenant_id = ? AND request_id = ?", tenant, approvedPermission.RequestID).Error).To(HaveOccurred())
	staleResult, err := executor.Execute(requesterCtx, approvedStale.RequestID, approvedStale.Version)
	Expect(err).NotTo(HaveOccurred())
	Expect(staleResult.Status).To(Equal(kind.Stale))
	_, err = executor.Execute(requesterCtx, approvedAssignment.RequestID, approvedAssignment.Version)
	Expect(err).To(MatchError(authority.ErrLastAdministrator))
	preservedAssignment, err := assignmentimpl.New(db).Find(ctx, tenant, requester)
	Expect(err).NotTo(HaveOccurred())
	Expect(preservedAssignment).To(Equal(originalAssignment))
	Expect(db.Raw("SELECT * FROM identity.approval_requests WHERE tenant_id = ? AND request_id = ?", tenant, approvedAssignment.RequestID).Scan(&persisted).Error).NotTo(HaveOccurred())
	Expect(persisted.Status).To(Equal(kind.Approved))
	roleApplied, err := executor.Execute(requesterCtx, approved.RequestID, approved.Version)
	Expect(err).NotTo(HaveOccurred())
	Expect(roleApplied.Status).To(Equal(kind.Applied))
	roleAfter, err := roleimpl.New(db).FindByID(ctx, tenant, targetRole)
	Expect(err).NotTo(HaveOccurred())
	Expect(roleAfter.Role.Status).To(Equal(entity.RoleInactive))
	Expect(db.Raw("SELECT count(*) FROM identity.resource_audit WHERE tenant_id = ? AND request_id = ? AND event_payload#>>'{data,old,role,status}' = 'ACTIVE' AND event_payload#>>'{data,new,role,status}' = 'INACTIVE'", tenant, approved.RequestID).Scan(&auditCount).Error).NotTo(HaveOccurred())
	Expect(auditCount).To(Equal(int64(1)))
	Expect(db.Raw("SELECT count(*) FROM identity.resource_audit WHERE tenant_id = ? AND request_id IN ?", tenant, []string{approvedAssignment.RequestID, approvedStale.RequestID}).Scan(&auditCount).Error).NotTo(HaveOccurred())
	Expect(auditCount).To(BeZero())

	const targetAccount = "10000000-0000-4000-8000-000000000015"
	Expect(db.Exec("INSERT INTO identity.accounts(account_id, person_gid, tenant_id, email, status) VALUES (?, ?, ?, 'resource-audit-target@approval.test', 'ACTIVE')", targetAccount, targetAccount, tenant).Error).NotTo(HaveOccurred())
	assignment.Intent.AccountID = targetAccount
	assignment.Intent.PrimaryRoleID = requiredRole
	assignment.Intent.RoleVersions = []accountdto.RoleVersion{{RoleID: requiredRole, Version: 1}}
	assignment.IdempotencyKey = "resource-audit-assignment"
	pendingTarget, err := service.CreateAssignment(requesterCtx, assignment)
	Expect(err).NotTo(HaveOccurred())
	approvedTarget, err := service.Decide(reviewerCtx, dto.Review{RequestID: pendingTarget.RequestID, Version: 1, Decision: kind.Approve})
	Expect(err).NotTo(HaveOccurred())
	appliedTarget, err := executor.Execute(requesterCtx, approvedTarget.RequestID, approvedTarget.Version)
	Expect(err).NotTo(HaveOccurred())
	Expect(appliedTarget.Status).To(Equal(kind.Applied))
	events, err = reads.Audit(requesterCtx, appliedTarget.RequestID)
	Expect(err).NotTo(HaveOccurred())
	Expect(events).To(HaveLen(4))
	Expect(db.Raw("SELECT count(*) FROM identity.resource_audit a JOIN infrastructure.messaging m ON m.event_id = a.event_id::text WHERE a.tenant_id = ? AND a.request_id = ? AND a.resource_id = ? AND a.event_payload = convert_from(m.payload, 'UTF8')::jsonb AND a.event_payload->>'eventType' = 'identity.account.authority.changed.v1' AND a.event_payload#>>'{data,old,account,primaryRoleId}' = '' AND a.event_payload#>>'{data,new,account,primaryRoleId}' = ?", tenant, approvedTarget.RequestID, targetAccount, requiredRole).Scan(&auditCount).Error).NotTo(HaveOccurred())
	Expect(auditCount).To(Equal(int64(1)))
}


var errInjectedEffect = errors.New("injected failure after draft write")

type failingDraftWriter struct {
	rolerepo.DraftPermissionWriter
}

func (writer failingDraftWriter) ReplaceDraftPermissions(ctx context.Context, request rolerepo.DraftPermissionUpdate) (rolerepo.PermissionSet, error) {
	if _, err := writer.DraftPermissionWriter.ReplaceDraftPermissions(ctx, request); err != nil {
		return rolerepo.PermissionSet{}, err
	}
	return rolerepo.PermissionSet{}, errInjectedEffect
}



var errInjectedOutbox = errors.New("injected approval outbox failure")

type failingOutbox struct{}

func (failingOutbox) Enqueue(ctx context.Context, event messaging.Event) error {
	return errInjectedOutbox
}

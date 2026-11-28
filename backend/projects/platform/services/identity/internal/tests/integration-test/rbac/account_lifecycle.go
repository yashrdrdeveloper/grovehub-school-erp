package rbac_test

import (
	"context"
	"grovehub/services/audit"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	schoolgateway "grovehub/launcher/session"
	kind "grovehub/services/identity/internal/app/constants/approval"
	"grovehub/services/identity/internal/app/constants/feature"
	"grovehub/services/identity/internal/app/constants/permission"
	approvalcontroller "grovehub/services/identity/internal/app/controller/approval/protected/impl"
	dto "grovehub/services/identity/internal/app/dto/approval"
	"grovehub/services/identity/internal/app/entity"
	accountproducer "grovehub/services/identity/internal/app/events/producers/account/impl"
	approvalproducer "grovehub/services/identity/internal/app/events/producers/approval"
	resourceproducer "grovehub/services/identity/internal/app/events/producers/resource/impl"
	mapper "grovehub/services/identity/internal/app/mappers/account/impl"
	accountimpl "grovehub/services/identity/internal/app/repository/account/impl"
	assignmentimpl "grovehub/services/identity/internal/app/repository/accountrole/impl"
	approvalrepo "grovehub/services/identity/internal/app/repository/approval/impl"
	policyimpl "grovehub/services/identity/internal/app/repository/approvalpolicy/impl"
	authorityimpl "grovehub/services/identity/internal/app/repository/authority/impl"
	authrepo "grovehub/services/identity/internal/app/repository/authorization/impl"
	registryrepo "grovehub/services/identity/internal/app/repository/registry/impl"
	roleimpl "grovehub/services/identity/internal/app/repository/role/impl"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
	sessionimpl "grovehub/services/identity/internal/app/repository/session/impl"
	"grovehub/services/identity/internal/app/services/account/operations/transition"
	"grovehub/services/identity/internal/app/services/approval"
	approvalimpl "grovehub/services/identity/internal/app/services/approval/impl"
	"grovehub/services/identity/internal/app/services/approval/operations/cutover"
	authimpl "grovehub/services/identity/internal/app/services/authorization/impl"
	"grovehub/services/identity/internal/app/services/authorization/operations/protect"
	"grovehub/services/identity/internal/app/services/registry"
	registryimpl "grovehub/services/identity/internal/app/services/registry/impl"
	messaging "grovehub/services/infrastructure/client/messaging"
)

func verifyAccountLifecycleCutover(ctx context.Context, db *gorm.DB, tenant string, requester string, reviewer string, managementRole string, target string, service approval.Service, store messaging.EventStore, recorder approvalproducer.Recorder) {
	codes := []permission.Code{permission.AccountSuspend, permission.AccountActivate, permission.AccountDeactivate}
	_, err := registryimpl.New(registryrepo.New(db)).Register(ctx, registry.Request{Features: []feature.Code{feature.AccountManagement}, Permissions: codes})
	Expect(err).NotTo(HaveOccurred())
	Expect(db.Exec("UPDATE identity.features SET active = true WHERE code = ?", feature.AccountManagement).Error).NotTo(HaveOccurred())
	Expect(db.Exec("UPDATE identity.permissions SET active = true WHERE code IN ?", codes).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.tenant_features(tenant_id, feature_id, enabled) SELECT ?, feature_id, true FROM identity.features WHERE code = ?", tenant, feature.AccountManagement).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.role_permissions(tenant_id, role_id, permission_id) SELECT ?, ?, permission_id FROM identity.permissions WHERE code IN ?", tenant, managementRole, codes).Error).NotTo(HaveOccurred())

	now := time.Now().UTC()
	requesterCtx := schoolgateway.WithPrincipal(ctx, schoolgateway.Principal{TenantID: tenant, AccountID: requester, SessionID: "requester-session", AuthenticatedAt: now})
	reviewerCtx := schoolgateway.WithPrincipal(ctx, schoolgateway.Principal{TenantID: tenant, AccountID: reviewer, SessionID: "reviewer-session", AuthenticatedAt: now})
	locker := authorityimpl.New(db)
	guard := protect.New(locker, authorityimpl.NewAdministratorReader(db))
	readiness := policyimpl.NewReadiness(db)
	accounts := accountimpl.NewNormalizedReader(db)
	assignments := assignmentimpl.New(db)
	snapshots := authrepo.New(db)
	authorizer := authimpl.New(snapshots)
	sessions := sessionimpl.New(db)
	effects := approvalimpl.Effects{
		Audit: resourceproducer.New(audit.NewProducer(store)),
		AccountLifecycle: transition.NewNormalized(guard, readiness, authorizer, accounts, accountimpl.NewStatusWriter(db), assignments, sessions, mapper.New(), accountproducer.New(store)),
	}
	build := func(effects approvalimpl.Effects) approval.Executor {
		return cutover.Wrap(cutover.New(tenant), approvalimpl.NewExecutor(approvalrepo.New(db, recorder), locker, guard, readiness, authorizer, snapshots, accounts, roleimpl.New(db), assignments, effects))
	}
	executor := build(effects)
	firstSession := entity.Session{SessionID: uuid.NewString(), AccountID: target, AccessHash: strings.Repeat("a",64), RefreshHash: strings.Repeat("b",64), AccessExpiresAt: now.Add(time.Hour), RefreshExpiresAt: now.Add(24*time.Hour), CreatedAt: now, UpdatedAt: now}
	Expect(sessions.Create(ctx, firstSession)).To(Succeed())

	submit := func(action string, version int64) entity.ApprovalRequest {
		pending, err := service.CreateLifecycle(requesterCtx, dto.CreateLifecycle{Type: kind.AccountLifecycleChange, Intent: dto.LifecycleIntent{TargetID: target, Version: version, Action: action}, Reason: "Lifecycle verification", IdempotencyKey: "cutover-"+action})
		Expect(err).NotTo(HaveOccurred())
		approved, err := service.Decide(reviewerCtx, dto.Review{RequestID: pending.RequestID, Version: pending.Version, Decision: kind.Approve})
		Expect(err).NotTo(HaveOccurred())
		Expect(approved.Status).To(Equal(kind.Approved))
		return approved
	}
	before, err := accounts.FindByID(ctx, target)
	Expect(err).NotTo(HaveOccurred())
	suspended := submit("suspend", before.Version)
	fault := effects
	fault.Audit = resourceproducer.New(audit.NewProducer(failingOutbox{}))
	_, err = build(fault).Execute(requesterCtx, suspended.RequestID, suspended.Version)
	Expect(err).To(MatchError(errInjectedOutbox))
	preserved, err := accounts.FindByID(ctx, target)
	Expect(err).NotTo(HaveOccurred())
	Expect(preserved.Status).To(Equal(entity.Active))
	Expect(preserved.Version).To(Equal(before.Version))
	_, err = sessions.FindAccess(ctx, firstSession.AccessHash, now)
	Expect(err).NotTo(HaveOccurred())

	controller := approvalcontroller.NewWithMutations(nil, service, executor)
	request := httptest.NewRequest(http.MethodPost, "/identity/protected/approvals/"+suspended.RequestID+"/execute", strings.NewReader(`{"version":2}`)).WithContext(requesterCtx)
	request.SetPathValue("requestId", suspended.RequestID)
	response := httptest.NewRecorder()
	controller.Execute(response, request)
	Expect(response.Code).To(Equal(http.StatusOK))
	Expect(response.Body.String()).To(ContainSubstring("APPLIED"))
	after, err := accounts.FindByID(ctx, target)
	Expect(err).NotTo(HaveOccurred())
	Expect(after.Status).To(Equal(entity.Suspended))
	Expect(after.Version).To(Equal(before.Version+1))
	_, err = sessions.FindAccess(ctx, firstSession.AccessHash, now)
	Expect(err).To(MatchError(sessionrepo.ErrInvalidSession))
	_, err = sessions.FindRefresh(ctx, firstSession.RefreshHash, now)
	Expect(err).To(MatchError(sessionrepo.ErrInvalidSession))
	secondSession := firstSession
	secondSession.SessionID = uuid.NewString()
	secondSession.AccessHash = strings.Repeat("c",64)
	secondSession.RefreshHash = strings.Repeat("d",64)
	Expect(sessions.Create(ctx, secondSession)).To(MatchError(sessionrepo.ErrInvalidSession))
	var count int64
	Expect(db.Raw("SELECT count(*) FROM (SELECT event_id::uuid AS event_id, p#>>'{data,tenantId}' AS tenant_id, (p#>>'{data,requestId}')::uuid AS request_id, (p->>'aggregateId')::uuid AS resource_id, (p#>>'{data,new,version}')::bigint AS resource_version, p#>>'{data,operation}' AS operation, p AS event_payload FROM (SELECT event_id, convert_from(payload,'UTF8')::jsonb AS p FROM infrastructure.messaging) x WHERE p->>'eventType' IN ('identity.role.changed.v1','identity.account.authority.changed.v1')) a WHERE tenant_id = ? AND request_id = ? AND event_payload#>>'{data,old,account,status}' = 'ACTIVE' AND event_payload#>>'{data,new,account,status}' = 'SUSPENDED'", tenant, suspended.RequestID).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(Equal(int64(1)))
	activated := submit("activate", after.Version)
	applied, err := executor.Execute(requesterCtx, activated.RequestID, activated.Version)
	Expect(err).NotTo(HaveOccurred())
	Expect(applied.Status).To(Equal(kind.Applied))
	_, err = sessions.FindAccess(ctx, firstSession.AccessHash, now)
	Expect(err).To(MatchError(sessionrepo.ErrInvalidSession))
	Expect(sessions.Create(ctx, secondSession)).To(Succeed())
	after, err = accounts.FindByID(ctx, target)
	Expect(err).NotTo(HaveOccurred())
	Expect(after.Status).To(Equal(entity.Active))
	deactivated := submit("deactivate", after.Version)
	applied, err = executor.Execute(requesterCtx, deactivated.RequestID, deactivated.Version)
	Expect(err).NotTo(HaveOccurred())
	Expect(applied.Status).To(Equal(kind.Applied))
	after, err = accounts.FindByID(ctx, target)
	Expect(err).NotTo(HaveOccurred())
	Expect(after.Status).To(Equal(entity.Deactivated))
	_, err = sessions.FindAccess(ctx, secondSession.AccessHash, now)
	Expect(err).To(MatchError(sessionrepo.ErrInvalidSession))
}

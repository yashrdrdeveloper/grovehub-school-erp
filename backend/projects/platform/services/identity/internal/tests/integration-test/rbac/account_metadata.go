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
	"grovehub/services/identity/internal/app/constants/permission"
	controller "grovehub/services/identity/internal/app/controller/account/protected/impl"
	dto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/entity"
	eventimpl "grovehub/services/identity/internal/app/events/producers/accountmetadata/impl"
	mapper "grovehub/services/identity/internal/app/mappers/account/impl"
	repo "grovehub/services/identity/internal/app/repository/account"
	accountimpl "grovehub/services/identity/internal/app/repository/account/impl"
	assignmentimpl "grovehub/services/identity/internal/app/repository/accountrole/impl"
	policyimpl "grovehub/services/identity/internal/app/repository/approvalpolicy/impl"
	authorityimpl "grovehub/services/identity/internal/app/repository/authority/impl"
	authrepo "grovehub/services/identity/internal/app/repository/authorization/impl"
	tokenrepo "grovehub/services/identity/internal/app/repository/passwordtoken"
	tokenimpl "grovehub/services/identity/internal/app/repository/passwordtoken/impl"
	registryrepo "grovehub/services/identity/internal/app/repository/registry/impl"
	accountservice "grovehub/services/identity/internal/app/services/account/impl"
	"grovehub/services/identity/internal/app/services/account/operations/update"
	updateimpl "grovehub/services/identity/internal/app/services/account/operations/update/impl"
	"grovehub/services/identity/internal/app/services/approval/operations/cutover"
	"grovehub/services/identity/internal/app/services/authorization"
	authimpl "grovehub/services/identity/internal/app/services/authorization/impl"
	"grovehub/services/identity/internal/app/services/registry"
	registryimpl "grovehub/services/identity/internal/app/services/registry/impl"
	validator "grovehub/services/identity/internal/app/validators/account/impl"
	messaging "grovehub/services/infrastructure/client/messaging"
)

func verifyAccountMetadata(ctx context.Context, db *gorm.DB, tenant string, requester string, reviewer string, managementRole string, target string, store messaging.EventStore) {
	_, err := registryimpl.New(registryrepo.New(db)).Register(ctx, registry.Request{Permissions: []permission.Code{permission.AccountUpdate}})
	Expect(err).NotTo(HaveOccurred())
	Expect(db.Exec("UPDATE identity.permissions SET active = true WHERE code = ?", permission.AccountUpdate).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.role_permissions(tenant_id, role_id, permission_id) SELECT ?, ?, permission_id FROM identity.permissions WHERE code = ? ON CONFLICT DO NOTHING", tenant, managementRole, permission.AccountUpdate).Error).NotTo(HaveOccurred())
	now := time.Now().UTC()
	requesterCtx := schoolgateway.WithPrincipal(ctx, schoolgateway.Principal{TenantID: tenant, AccountID: requester, SessionID: "requester-session", AuthenticatedAt: now})
	reviewerCtx := schoolgateway.WithPrincipal(ctx, schoolgateway.Principal{TenantID: tenant, AccountID: reviewer, SessionID: "reviewer-session", AuthenticatedAt: now})
	accounts := accountimpl.NewNormalizedReader(db)
	tokens := tokenimpl.New(db)
	locker := authorityimpl.New(db)
	authorizer := authimpl.New(authrepo.New(db))
	build := func(events messaging.EventStore) update.Operation {
		return updateimpl.NewNormalized(accounts, accountimpl.NewMetadataWriter(db), assignmentimpl.New(db), tokens, locker, policyimpl.NewReadiness(db), authorizer, mapper.New(), validator.NewNormalized(), eventimpl.New(audit.NewProducer(events)))
	}
	operation := build(store)
	before, err := accounts.FindByID(ctx, target)
	Expect(err).NotTo(HaveOccurred())
	username := " Updated.Account "
	request := dto.Update{Username: &username, Version: &before.Version}
	_, err = operation.Execute(reviewerCtx, target, request)
	Expect(err).To(MatchError(authorization.ErrForbidden))
	_, err = operation.Execute(requesterCtx, target, dto.Update{Username: &username})
	Expect(err).To(MatchError(repo.ErrConflict))
	_, err = operation.Execute(requesterCtx, target, dto.Update{Username: &username, Version: &before.Version, RoleAssignment: &dto.RoleAssignment{PrimaryRoleID: managementRole}})
	Expect(err).To(MatchError(repo.ErrConflict))
	_, err = build(failingOutbox{}).Execute(requesterCtx, target, request)
	Expect(err).To(MatchError(errInjectedOutbox))
	preserved, err := accounts.FindByID(ctx, target)
	Expect(err).NotTo(HaveOccurred())
	Expect(preserved).To(Equal(before))

	service := accountservice.NewOnboarding(nil, cutover.New(tenant), nil, nil, nil, nil, operation)
	api := controller.New(service)
	httpUpdate := func(body string) *httptest.ResponseRecorder {
		request := httptest.NewRequest(http.MethodPatch, "/identity/protected/accounts/"+target, strings.NewReader(body)).WithContext(requesterCtx)
		request.SetPathValue("accountId", target)
		response := httptest.NewRecorder()
		api.Update(response, request)
		return response
	}
	Expect(httpUpdate(`{"email":"unverified@example.com","version":3}`).Code).To(Equal(http.StatusBadRequest))
	Expect(httpUpdate(`{"phoneNumber":"9999999999","version":3}`).Code).To(Equal(http.StatusBadRequest))
	updated, err := service.Update(requesterCtx, target, request)
	Expect(err).NotTo(HaveOccurred())
	Expect(*updated.Username).To(Equal("updated.account"))
	Expect(updated.Version).To(Equal(before.Version+1))
	Expect(updated.PrimaryRoleID).To(Equal(before.PrimaryRoleID))
	Expect(updated.Email).To(Equal(before.Email))
	Expect(updated.Status).To(Equal(string(before.Status)))
	_, err = operation.Execute(requesterCtx, target, request)
	Expect(err).To(MatchError(repo.ErrVersionConflict))
	request.Version = &updated.Version
	replayed, err := operation.Execute(requesterCtx, target, request)
	Expect(err).NotTo(HaveOccurred())
	Expect(replayed.Version).To(Equal(updated.Version))
	var count int64
	Expect(db.Raw("SELECT count(*) FROM (SELECT event_id::uuid AS event_id, p#>>'{data,tenantId}' AS tenant_id, (p->>'aggregateId')::uuid AS account_id, p#>>'{data,operation}' AS operation, p AS event_payload FROM (SELECT event_id, convert_from(payload,'UTF8')::jsonb AS p FROM infrastructure.messaging) x WHERE p->>'eventType' = 'identity.account.metadata.changed.v1') a WHERE a.tenant_id = ? AND a.account_id = ? AND a.event_payload#>>'{data,new,username}' = 'updated.account'", tenant, target).Scan(&count).Error).NotTo(HaveOccurred())
	Expect(count).To(Equal(int64(1)))

	pending := uuid.NewString()
	Expect(db.Exec("INSERT INTO identity.accounts(account_id,person_gid,tenant_id,email,status) VALUES (?,?,?,?, 'PENDING_SETUP')", pending, uuid.NewString(), tenant, "metadata-pending@approval.test").Error).NotTo(HaveOccurred())
	version := int64(1)
	proofHash := strings.Repeat("e",64)
	grantHash := strings.Repeat("f",64)
	for _, proof := range []entity.PasswordToken{
		{TokenHash: proofHash, AccountID: pending, Purpose: entity.AccountSetup, AccountVersion: &version, CreatedAt: now, ExpiresAt: now.Add(time.Hour)},
		{TokenHash: grantHash, AccountID: pending, Purpose: entity.SetupGrant, AccountVersion: &version, CreatedAt: now, ExpiresAt: now.Add(time.Hour)},
	} {
		Expect(tokens.Issue(ctx, proof)).To(Succeed())
	}
	_, err = operation.Execute(requesterCtx, pending, dto.Update{Username: &username, Version: &version})
	Expect(err).To(MatchError(repo.ErrConflict))
	_, err = tokens.Inspect(ctx, grantHash, entity.SetupGrant, now)
	Expect(err).NotTo(HaveOccurred())
	username = "pending.metadata"
	_, err = build(failingOutbox{}).Execute(requesterCtx, pending, dto.Update{Username: &username, Version: &version})
	Expect(err).To(MatchError(errInjectedOutbox))
	_, err = tokens.Inspect(ctx, grantHash, entity.SetupGrant, now)
	Expect(err).NotTo(HaveOccurred())
	_, err = operation.Execute(requesterCtx, pending, dto.Update{Username: &username, Version: &version})
	Expect(err).NotTo(HaveOccurred())
	_, err = tokens.Inspect(ctx, proofHash, entity.AccountSetup, now)
	Expect(err).To(MatchError(tokenrepo.ErrInvalidToken))
	_, err = tokens.Inspect(ctx, grantHash, entity.SetupGrant, now)
	Expect(err).To(MatchError(tokenrepo.ErrInvalidToken))
}

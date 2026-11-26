package rbac_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/google/uuid"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"grovehub/libs/schoolgateway"
	"grovehub/services/identity/internal/app/constants/feature"
	"grovehub/services/identity/internal/app/constants/permission"
	controller "grovehub/services/identity/internal/app/controller/audit/protected/impl"
	dto "grovehub/services/identity/internal/app/dto/audit"
	auditrepo "grovehub/services/identity/internal/app/repository/auditread/impl"
	authorityrepo "grovehub/services/identity/internal/app/repository/authority/impl"
	authrepo "grovehub/services/identity/internal/app/repository/authorization/impl"
	registryrepo "grovehub/services/identity/internal/app/repository/registry/impl"
	"grovehub/services/identity/internal/app/services/audit"
	auditimpl "grovehub/services/identity/internal/app/services/audit/impl"
	"grovehub/services/identity/internal/app/services/authorization"
	authimpl "grovehub/services/identity/internal/app/services/authorization/impl"
	"grovehub/services/identity/internal/app/services/registry"
	registryimpl "grovehub/services/identity/internal/app/services/registry/impl"
)

func verifyAuditSearch(ctx context.Context, db *gorm.DB, tenant string, requester string, reviewer string, managementRole string) {
	authorizer := authimpl.New(authrepo.New(db))
	locker := authorityrepo.New(db)
	reader := auditrepo.New(db)
	service := auditimpl.New(reader, locker, authorizer, false)
	requesterCtx := schoolgateway.WithPrincipal(ctx, schoolgateway.Principal{TenantID: tenant, AccountID: requester, SessionID: "requester-session"})
	reviewerCtx := schoolgateway.WithPrincipal(ctx, schoolgateway.Principal{TenantID: tenant, AccountID: reviewer, SessionID: "reviewer-session"})
	_, err := service.Search(requesterCtx, dto.Search{})
	Expect(err).To(MatchError(authorization.ErrForbidden))
	_, err = registryimpl.New(registryrepo.New(db)).Register(ctx, registry.Request{Features: []feature.Code{feature.AuditManagement}, Permissions: []permission.Code{permission.AuditView}})
	Expect(err).NotTo(HaveOccurred())
	Expect(db.Exec("UPDATE identity.features SET active = true WHERE code = ?", feature.AuditManagement).Error).NotTo(HaveOccurred())
	Expect(db.Exec("UPDATE identity.permissions SET active = true WHERE code = ?", permission.AuditView).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.tenant_features(tenant_id, feature_id, enabled) SELECT ?, feature_id, true FROM identity.features WHERE code = ?", tenant, feature.AuditManagement).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.role_permissions(tenant_id, role_id, permission_id) SELECT ?, ?, permission_id FROM identity.permissions WHERE code = ?", tenant, managementRole, permission.AuditView).Error).NotTo(HaveOccurred())
	_, err = service.Search(reviewerCtx, dto.Search{})
	Expect(err).To(MatchError(authorization.ErrForbidden))
	_, err = auditimpl.New(reader, locker, authorizer, true).Search(requesterCtx, dto.Search{})
	Expect(err).To(MatchError(authorization.ErrForbidden))
	_, err = service.Search(requesterCtx, dto.Search{Limit: 101})
	Expect(err).To(MatchError(audit.ErrInvalidSearch))
	first, err := service.Search(requesterCtx, dto.Search{Limit: 1})
	Expect(err).NotTo(HaveOccurred())
	Expect(first.Items).To(HaveLen(1))
	Expect(first.NextCursor).NotTo(BeNil())
	second, err := service.Search(requesterCtx, dto.Search{Limit: 1, Cursor: first.NextCursor})
	Expect(err).NotTo(HaveOccurred())
	Expect(second.Items).To(HaveLen(1))
	Expect(second.Items[0].EventID).NotTo(Equal(first.Items[0].EventID))
	for _, source := range []string{"APPROVAL","RESOURCE","ROLE_METADATA","ACCOUNT_METADATA","REGISTRATION"} {
		page, err := service.Search(requesterCtx, dto.Search{Source: source})
		Expect(err).NotTo(HaveOccurred())
		Expect(page.Items).NotTo(BeEmpty())
		for _, event := range page.Items { Expect(event.Source).To(Equal(source)) }
	}
	page, err := service.Search(requesterCtx, dto.Search{ResourceID: uuid.NewString()})
	Expect(err).NotTo(HaveOccurred())
	Expect(page.Items).To(BeEmpty())
	foreign, err := reader.Search(ctx, "unknown-audit-tenant", dto.Search{Limit: 25})
	Expect(err).NotTo(HaveOccurred())
	Expect(foreign).To(BeEmpty())

	api := controller.New(service)
	request := httptest.NewRequest(http.MethodPost, "/identity/protected/audit/search", strings.NewReader(`{"source":"REGISTRATION"}`)).WithContext(requesterCtx)
	response := httptest.NewRecorder()
	api.Search(response, request)
	Expect(response.Code).To(Equal(http.StatusOK))
	Expect(response.Header().Get("Cache-Control")).To(Equal("no-store"))
	Expect(response.Body.String()).NotTo(ContainSubstring("eventPayload"))
	Expect(response.Body.String()).NotTo(ContainSubstring("requestSha256"))
	Expect(db.Exec("UPDATE identity.permissions SET active = false WHERE code = ?", permission.AuditView).Error).NotTo(HaveOccurred())
	_, err = service.Search(requesterCtx, dto.Search{})
	Expect(err).To(MatchError(authorization.ErrForbidden))
}

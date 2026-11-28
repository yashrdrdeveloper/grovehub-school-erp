package rbac_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	schoolgateway "grovehub/libs/platform/session"
	"grovehub/services/identity/internal/app/constants/permission"
	controller "grovehub/services/identity/internal/app/controller/capability/protected/impl"
	authorityimpl "grovehub/services/identity/internal/app/repository/authority/impl"
	authrepo "grovehub/services/identity/internal/app/repository/authorization/impl"
	capabilityrepo "grovehub/services/identity/internal/app/repository/capability/impl"
	registryrepo "grovehub/services/identity/internal/app/repository/registry/impl"
	authimpl "grovehub/services/identity/internal/app/services/authorization/impl"
	"grovehub/services/identity/internal/app/services/authorization"
	capabilityimpl "grovehub/services/identity/internal/app/services/capability/impl"
	"grovehub/services/identity/internal/app/services/registry"
	registryimpl "grovehub/services/identity/internal/app/services/registry/impl"
	"grovehub/services/identity/internal/app/services/security/operations/routepolicy"
)

func verifyCapabilityDiscovery(ctx context.Context, db *gorm.DB, tenant string, requester string, reviewer string, managementRole string) {
	_, err := registryimpl.New(registryrepo.New(db)).Register(ctx, registry.Request{Permissions: []permission.Code{permission.PermissionView, permission.FeatureManage}})
	Expect(err).NotTo(HaveOccurred())
	Expect(db.Exec("UPDATE identity.permissions SET active = true WHERE code IN ?", []permission.Code{permission.PermissionView, permission.FeatureManage}).Error).NotTo(HaveOccurred())
	Expect(db.Exec("INSERT INTO identity.role_permissions(tenant_id, role_id, permission_id) SELECT ?, ?, permission_id FROM identity.permissions WHERE code = ?", tenant, managementRole, permission.PermissionView).Error).NotTo(HaveOccurred())
	now := time.Now().UTC()
	actor := schoolgateway.Principal{TenantID: tenant, AccountID: requester, SessionID: "requester-session", AuthenticatedAt: now}
	requesterCtx := schoolgateway.WithPrincipal(ctx, actor)
	reviewerCtx := schoolgateway.WithPrincipal(ctx, schoolgateway.Principal{TenantID: tenant, AccountID: reviewer, SessionID: "reviewer-session"})
	authorizer := authimpl.New(authrepo.New(db))
	service := capabilityimpl.New(capabilityrepo.New(db), authorityimpl.New(db), authorizer, tenant)
	discovery, err := service.Discover(requesterCtx)
	Expect(err).NotTo(HaveOccurred())
	Expect(discovery.Permissions).To(ContainElement(permission.RolePermissionManage))
	Expect(discovery.Actions["rolePermissionsSubmit"]).To(BeTrue())
	Expect(discovery.Actions["accountCreate"]).To(BeFalse())
	Expect(discovery.Actions["roleCreate"]).To(BeFalse())
	Expect(discovery.Actions["permissionCatalog"]).To(BeTrue())
	Expect(discovery.ExecutionConfigured).To(BeTrue())
	Expect(discovery.ExecutionRequiresReauthentication).To(BeFalse())
	actor.AuthenticatedAt = now.Add(-10*time.Minute)
	discovery, err = service.Discover(schoolgateway.WithPrincipal(ctx, actor))
	Expect(err).NotTo(HaveOccurred())
	Expect(discovery.ExecutionRequiresReauthentication).To(BeTrue())

	catalog, err := service.Catalog(requesterCtx)
	Expect(err).NotTo(HaveOccurred())
	found := false
	for _, item := range catalog.Permissions {
		Expect(item.Code).NotTo(Equal(permission.FeatureManage))
		if item.Code == permission.RoleDeactivate {
			found = true
			Expect(item.CanDelegate).To(BeTrue())
		}
	}
	Expect(found).To(BeTrue())
	_, err = service.Catalog(reviewerCtx)
	Expect(err).To(MatchError(authorization.ErrForbidden))
	_, err = service.Discover(schoolgateway.WithPrincipal(ctx, schoolgateway.Principal{TenantID: "foreign", AccountID: requester, SessionID: "foreign-session"}))
	Expect(err).To(MatchError(authorization.ErrForbidden))
	api := controller.New(service)
	response := httptest.NewRecorder()
	api.Discover(response, httptest.NewRequest(http.MethodGet, "/identity/protected/capabilities/me", nil).WithContext(requesterCtx))
	Expect(response.Code).To(Equal(http.StatusOK))
	Expect(response.Header().Get("Cache-Control")).To(Equal("no-store"))
	Expect(response.Body.String()).To(ContainSubstring("executionConfigured"))
	policy := routepolicy.New(authorizer, false)
	Expect(policy.Authorize(requesterCtx, "GET /identity/protected/capabilities/me")).To(Succeed())
	Expect(policy.Authorize(requesterCtx, "GET /identity/protected/capabilities/permissions")).To(Succeed())
	Expect(policy.Authorize(reviewerCtx, "GET /identity/protected/capabilities/permissions")).To(MatchError(authorization.ErrForbidden))
	Expect(routepolicy.New(authorizer, true).Authorize(requesterCtx, "GET /identity/protected/capabilities/me")).To(MatchError(authorization.ErrForbidden))
}

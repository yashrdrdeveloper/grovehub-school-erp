package routepolicy_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/schoolgateway"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/security/operations/routepolicy"
)

type authorizer struct {
	code permission.Code
	err error
}

func (auth *authorizer) Require(ctx context.Context, code permission.Code) error {
	auth.code = code
	return auth.err
}

func (auth *authorizer) Effective(ctx context.Context) ([]permission.Code, error) {
	return nil, auth.err
}

func TestRoutePolicy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Identity Route Policy Suite")
}

var _ = Describe("Authorize", func() {
	It("requires live permissions for management reads and denies unknown and unsafe writes", func() {
		ctx := schoolgateway.WithPrincipal(context.Background(), schoolgateway.Principal{AccountID: "actor", TenantID: "school", SessionID: "session", Admin: true})
		auth := &authorizer{err: authorization.ErrForbidden}
		policy := routepolicy.New(auth, false)
		Expect(policy.Authorize(ctx, "GET /identity/protected/accounts/{accountId}")).To(MatchError(authorization.ErrForbidden))
		Expect(auth.code).To(Equal(permission.AccountView))
		auth.err = nil
		Expect(policy.Authorize(ctx, "POST /identity/protected/accounts/search")).To(Succeed())
		Expect(auth.code).To(Equal(permission.AccountSearch))
		Expect(policy.Authorize(ctx, "PATCH /identity/protected/accounts/{accountId}/suspend")).To(MatchError(authorization.ErrForbidden))
		Expect(policy.Authorize(ctx, "GET /identity/protected/unknown")).To(MatchError(authorization.ErrForbidden))
		Expect(policy.Authorize(ctx, "GET /identity/protected/auth/me")).To(Succeed())
		Expect(routepolicy.New(nil, false).Authorize(ctx, "POST /identity/protected/accounts/search")).To(MatchError(authorization.ErrUnavailable))
	})
	It("requires authentication even for self-service and explicit development compatibility", func() {
		policy := routepolicy.New(nil, true)
		Expect(policy.Authorize(context.Background(), "GET /identity/protected/auth/me")).NotTo(Succeed())
		Expect(policy.Authorize(context.Background(), "POST /identity/protected/accounts")).NotTo(Succeed())
		ctx := schoolgateway.WithPrincipal(context.Background(), schoolgateway.Principal{AccountID: "actor", TenantID: "school", SessionID: "session"})
		Expect(policy.Authorize(ctx, "POST /identity/protected/accounts")).To(Succeed())
		Expect(policy.Authorize(ctx, "POST /identity/protected/unknown")).To(MatchError(authorization.ErrForbidden))
	})
})

var _ = Describe("RBAC read routes", func() {
	It("requires the exact capability for role and assignment reads and previews", func() {
		ctx := schoolgateway.WithPrincipal(context.Background(), schoolgateway.Principal{AccountID: "actor", TenantID: "school", SessionID: "session", Admin: true})
		auth := &authorizer{}
		policy := routepolicy.New(auth, false)
		for pattern, code := range map[string]permission.Code{
			"POST /identity/protected/roles/search": permission.RoleView,
			"GET /identity/protected/roles/{roleId}": permission.RoleView,
			"GET /identity/protected/roles/{roleId}/permissions": permission.RolePermissionView,
			"POST /identity/protected/roles/{roleId}/permissions/preview": permission.RolePermissionManage,
			"GET /identity/protected/account-roles/{accountId}": permission.AccountRoleView,
			"POST /identity/protected/account-roles/{accountId}/preview": permission.AccountRoleAssign,
		} {
			Expect(policy.Authorize(ctx, pattern)).To(Succeed())
			Expect(auth.code).To(Equal(code))
			Expect(routepolicy.New(auth, true).Authorize(ctx, pattern)).To(MatchError(authorization.ErrForbidden))
		}
		Expect(policy.Authorize(ctx, "POST /identity/protected/roles")).To(MatchError(authorization.ErrForbidden))
	})
	It("requires tenant authentication and service-level visibility for approval reads", func() {
		ctx := schoolgateway.WithPrincipal(context.Background(), schoolgateway.Principal{AccountID: "actor", TenantID: "school", SessionID: "session"})
		for _, pattern := range []string{
			"POST /identity/protected/approvals/history/search",
			"POST /identity/protected/approvals/inbox/search",
			"GET /identity/protected/approvals/inbox/count",
			"GET /identity/protected/approvals/{requestId}",
			"GET /identity/protected/approvals/{requestId}/audit",
		} {
			Expect(routepolicy.New(nil, false).Authorize(ctx, pattern)).To(Succeed())
			Expect(routepolicy.New(nil, false).Authorize(context.Background(), pattern)).NotTo(Succeed())
			Expect(routepolicy.New(nil, true).Authorize(ctx, pattern)).To(MatchError(authorization.ErrForbidden))
		}
		Expect(routepolicy.New(nil, false).Authorize(ctx, "POST /identity/protected/approvals/{requestId}/execute")).To(MatchError(authorization.ErrForbidden))
	})
})

var _ = Describe("Approval mutation routes", func() {
	It("allows strict authenticated request-specific operations, but keeps execution blocked", func() {
		ctx := schoolgateway.WithPrincipal(context.Background(), schoolgateway.Principal{AccountID: "actor", TenantID: "school", SessionID: "session"})
		auth := &authorizer{}
		policy := routepolicy.New(auth, false)
		for _, path := range []string{"/lifecycle", "/{requestId}/approve", "/{requestId}/reject", "/{requestId}/cancel"} {
			pattern := "POST /identity/protected/approvals"+path
			Expect(policy.Authorize(ctx, pattern)).To(Succeed())
			Expect(policy.Authorize(context.Background(), pattern)).NotTo(Succeed())
			Expect(routepolicy.New(auth, true).Authorize(ctx, pattern)).To(MatchError(authorization.ErrForbidden))
		}
		Expect(policy.Authorize(ctx, "POST /identity/protected/approvals/permissions")).To(Succeed())
		Expect(auth.code).To(Equal(permission.RolePermissionManage))
		Expect(policy.Authorize(ctx, "POST /identity/protected/approvals/assignments")).To(Succeed())
		Expect(auth.code).To(Equal(permission.AccountRoleAssign))
		Expect(policy.Authorize(ctx, "POST /identity/protected/approvals/{requestId}/execute")).To(MatchError(authorization.ErrForbidden))
		Expect(routepolicy.New(auth, true).Authorize(ctx, "POST /identity/protected/approvals/{requestId}/execute")).To(MatchError(authorization.ErrForbidden))
	})
})

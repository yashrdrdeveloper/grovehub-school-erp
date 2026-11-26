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

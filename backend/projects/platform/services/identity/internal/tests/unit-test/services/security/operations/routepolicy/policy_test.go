package routepolicy_test

import (
	"context"
	"testing"
	"time"
	"grovehub/services/identity/internal/app/services/approval/operations/cutover"
	"grovehub/services/identity/internal/app/services/security/principal"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	schoolgateway "grovehub/launcher/services/gateway/session"
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
			"AUDIT tenant search": permission.AuditView,
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
		Expect(policy.Authorize(ctx, "POST /identity/protected/roles")).To(Succeed())
		Expect(auth.code).To(Equal(permission.RoleCreate))
		Expect(policy.Authorize(ctx, "PATCH /identity/protected/roles/{roleId}")).To(Succeed())
		Expect(auth.code).To(Equal(permission.RoleUpdate))
		Expect(routepolicy.New(auth, true).Authorize(ctx, "POST /identity/protected/roles")).To(MatchError(authorization.ErrForbidden))
		Expect(routepolicy.New(auth, true).Authorize(ctx, "PATCH /identity/protected/roles/{roleId}")).To(MatchError(authorization.ErrForbidden))
		Expect(policy.Authorize(ctx, "DELETE /identity/protected/roles/{roleId}")).To(MatchError(authorization.ErrForbidden))
	})
	It("requires tenant authentication and service-level visibility for approval reads", func() {
		ctx := schoolgateway.WithPrincipal(context.Background(), schoolgateway.Principal{AccountID: "actor", TenantID: "school", SessionID: "session"})
		for _, pattern := range []string{
			"POST /audit/protected/search",
			"POST /identity/protected/approvals/history/search",
			"POST /identity/protected/approvals/inbox/search",
			"GET /identity/protected/approvals/inbox/count",
			"GET /identity/protected/approvals/{requestId}",
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

var _ = Describe("Execution cutover", func() {
	It("allows only the explicitly configured tenant with recent authentication", func() {
		pattern := "POST /identity/protected/approvals/{requestId}/execute"
		now := time.Now().UTC()
		actor := schoolgateway.Principal{AccountID: "actor", TenantID: "school", SessionID: "session", AuthenticatedAt: now}
		ctx := schoolgateway.WithPrincipal(context.Background(), actor)
		Expect(routepolicy.NewWithExecution(nil, false, cutover.New("school")).Authorize(ctx, pattern)).To(Succeed())
		Expect(routepolicy.NewWithExecution(nil, false, cutover.New("")).Authorize(ctx, pattern)).To(MatchError(authorization.ErrForbidden))
		Expect(routepolicy.NewWithExecution(nil, false, cutover.New("other")).Authorize(ctx, pattern)).To(MatchError(authorization.ErrForbidden))
		Expect(routepolicy.NewWithExecution(nil, true, cutover.New("school")).Authorize(ctx, pattern)).To(MatchError(authorization.ErrForbidden))
		actor.AuthenticatedAt = now.Add(-10*time.Minute)
		Expect(routepolicy.NewWithExecution(nil, false, cutover.New("school")).Authorize(schoolgateway.WithPrincipal(context.Background(), actor), pattern)).To(MatchError(principal.ErrRecentAuthenticationRequired))
		Expect(routepolicy.NewWithExecution(nil, false, cutover.New("school")).Authorize(ctx, "PATCH /identity/protected/accounts/{accountId}/suspend")).To(MatchError(authorization.ErrForbidden))
	})
})

var _ = Describe("Onboarding cutover", func() {
 It("requires exact tenant, recent authentication and the operation permission", func() {
  actor := schoolgateway.Principal{AccountID: "actor", TenantID: "school", SessionID: "session", AuthenticatedAt: time.Now().UTC()}
  ctx := schoolgateway.WithPrincipal(context.Background(), actor)
  auth := &authorizer{}
  for pattern, code := range map[string]permission.Code{
   "POST /identity/protected/accounts": permission.AccountCreate,
   "POST /identity/protected/accounts/{accountId}/resend-invitation": permission.AccountResendInvitation,
   "PATCH /identity/protected/accounts/{accountId}": permission.AccountUpdate,
  } {
   Expect(routepolicy.NewWithExecution(auth, false, cutover.New("school")).Authorize(ctx, pattern)).To(Succeed())
   Expect(auth.code).To(Equal(code))
   Expect(routepolicy.NewWithExecution(auth, false, cutover.New("other")).Authorize(ctx, pattern)).To(MatchError(authorization.ErrForbidden))
   Expect(routepolicy.New(auth, false).Authorize(ctx, pattern)).To(MatchError(authorization.ErrForbidden))
  }
 })
})

func TestConfiguredTenantPolicyForDirectCallers(t *testing.T) {
 active, available := true, true
 generation := int64(2)
 calls := 0
 policy := routepolicy.NewWithExecution(&authorizer{}, false, cutover.New("school"), func(_ context.Context, tenant string)(int64,error){
  calls++
  if tenant!="school" || !active || !available { return 0,schoolgateway.ErrTenantDenied }
  return generation,nil
 })
 actor := schoolgateway.Principal{AccountID:"actor",TenantID:"school",SessionID:"session",SessionGeneration:2,AuthenticatedAt:time.Now().UTC()}
 ctx := schoolgateway.WithPrincipal(context.Background(),actor)
 paths := []string{
  "GET /identity/protected/auth/me",
  "GET /identity/protected/auth/sessions",
  "POST /identity/protected/auth/change-password",
  "GET /identity/protected/accounts/{accountId}",
  "GET /identity/protected/accounts/{accountId}/person",
  "POST /identity/protected/accounts/search",
  "POST /identity/protected/roles/search",
  "POST /identity/protected/approvals/{requestId}/execute",
  "POST /audit/protected/search",
  "AUDIT tenant search",
 }
 for _,path:=range paths{if err:=policy.Authorize(ctx,path);err!=nil{t.Fatal(path,err)}}
 active=false
 for _,path:=range paths{if err:=policy.Authorize(ctx,path);err!=authorization.ErrForbidden{t.Fatal("disabled tenant accepted",path,err)}}
 active=true;generation++
 for _,path:=range paths{if err:=policy.Authorize(ctx,path);err!=authorization.ErrForbidden{t.Fatal("stale session accepted",path,err)}}
 // A caller cannot override the configured policy with a permissive context.
 forged:=schoolgateway.WithTenantCheck(ctx,func(context.Context,string)(int64,error){return 2,nil})
 if err:=policy.Authorize(forged,paths[0]);err!=authorization.ErrForbidden{t.Fatal("caller policy overrode configuration",err)}
 actor.SessionGeneration=3
 ctx=schoolgateway.WithPrincipal(context.Background(),actor)
 available=false
 if err:=policy.Authorize(ctx,paths[0]);err!=authorization.ErrForbidden{t.Fatal("unavailable registry accepted",err)}
 before:=calls
 for _,path:=range []string{"POST /identity/protected/auth/logout","POST /identity/protected/auth/logout-all","DELETE /identity/protected/auth/sessions/{sessionId}"}{
  if err:=policy.Authorize(ctx,path);err!=nil{t.Fatal("cleanup denied",err)}
  if err:=policy.Authorize(context.Background(),path);err==nil{t.Fatal("anonymous cleanup accepted")}
 }
 if calls!=before{t.Fatal("cleanup depends on registry")}
}

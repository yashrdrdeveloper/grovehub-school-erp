package authorization_test

import (
	"context"
	"errors"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	schoolgateway "grovehub/internal/session"
	"grovehub/services/identity/internal/app/constants/feature"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/entity"
	authorizationrepo "grovehub/services/identity/internal/app/repository/authorization"
	"grovehub/services/identity/internal/app/services/authorization"
	authorizationimpl "grovehub/services/identity/internal/app/services/authorization/impl"
	"grovehub/services/identity/internal/app/services/security/principal"
)

func TestAuthorizationService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Authorization Service Suite")
}

type loader struct {
	snapshot authorizationrepo.Snapshot
	err      error
	calls    int
	tenant   string
	account  string
}

func (repository *loader) LoadSnapshot(ctx context.Context, tenantID string, accountID string) (authorizationrepo.Snapshot, error) {
	repository.calls++
	repository.tenant = tenantID
	repository.account = accountID
	return repository.snapshot, repository.err
}

var _ = Describe("Authorization service", func() {
	var repo *loader
	var service authorization.Service
	var ctx context.Context

	BeforeEach(func() {
		repo = &loader{snapshot: authorizationrepo.Snapshot{
			AccountID: "actor", TenantID: "school", Status: entity.Active,
			Roles: []authorizationrepo.RoleGrant{
				{RoleID: "role", TenantID: "school", Active: true, Permissions: []permission.Code{permission.AccountSearch, permission.AccountSelfView}},
				{RoleID: "second", TenantID: "school", Active: true, Permissions: []permission.Code{permission.AccountSelfView}},
			},
			Permissions: map[permission.Code]authorizationrepo.PermissionState{
				permission.AccountSearch:   {Active: true, Feature: feature.AccountManagement},
				permission.AccountSelfView: {Active: true},
			},
			Features: map[feature.Code]authorizationrepo.FeatureState{feature.AccountManagement: {Active: true, Enabled: true}},
		}}
		ctx = schoolgateway.WithPrincipal(context.Background(), schoolgateway.Principal{AccountID: "actor", TenantID: "school", SessionID: "session", Admin: true})
		service = authorizationimpl.New(repo)
	})

	It("uses only the verified actor scope and rechecks authority on every call", func() {
		Expect(service.Require(ctx, permission.AccountSearch)).To(Succeed())
		Expect(repo.tenant).To(Equal("school"))
		Expect(repo.account).To(Equal("actor"))
		repo.snapshot.Features[feature.AccountManagement] = authorizationrepo.FeatureState{Active: true}
		Expect(service.Require(ctx, permission.AccountSearch)).To(MatchError(authorization.ErrForbidden))
		Expect(repo.calls).To(Equal(2))
	})

	It("lists effective permissions once with no duplicates or implicit admin access", func() {
		actual, err := service.Effective(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(actual).To(Equal([]permission.Code{permission.AccountSearch, permission.AccountSelfView}))
		Expect(repo.calls).To(Equal(1))
		Expect(service.Require(ctx, permission.AccountDeactivate)).To(MatchError(authorization.ErrForbidden))
	})

	It("rejects unauthenticated, tenantless and unknown permission requests before I/O", func() {
		Expect(service.Require(context.Background(), permission.AccountSearch)).To(MatchError(principal.ErrAuthenticationRequired))
		tenantless := schoolgateway.WithPrincipal(context.Background(), schoolgateway.Principal{AccountID: "actor", SessionID: "session"})
		Expect(service.Require(tenantless, permission.AccountSearch)).To(MatchError(principal.ErrTenantForbidden))
		Expect(service.Require(ctx, permission.Code("UNREGISTERED"))).To(MatchError(authorization.ErrForbidden))
		Expect(repo.calls).To(BeZero())
	})

	It("fails closed for missing, inactive or out-of-scope accounts", func() {
		repo.err = authorizationrepo.ErrNotFound
		Expect(service.Require(ctx, permission.AccountSearch)).To(MatchError(authorization.ErrForbidden))
		repo.err = nil
		repo.snapshot.Status = entity.Suspended
		Expect(service.Require(ctx, permission.AccountSearch)).To(MatchError(authorization.ErrForbidden))
		repo.snapshot.Status = entity.Active
		repo.snapshot.TenantID = "foreign"
		Expect(service.Require(ctx, permission.AccountSearch)).To(MatchError(authorization.ErrForbidden))
		repo.snapshot.TenantID = "school"
		repo.snapshot.AccountID = "foreign"
		actual, err := service.Effective(ctx)
		Expect(err).To(MatchError(authorization.ErrForbidden))
		Expect(actual).To(BeNil())
	})

	It("does not return grants on a repository failure", func() {
		repo.err = errors.New("database unavailable")
		err := service.Require(ctx, permission.AccountSearch)
		Expect(errors.Is(err, authorization.ErrUnavailable)).To(BeTrue())
		Expect(errors.Is(err, repo.err)).To(BeTrue())
		actual, err := service.Effective(ctx)
		Expect(errors.Is(err, authorization.ErrUnavailable)).To(BeTrue())
		Expect(actual).To(BeNil())
		Expect(authorizationimpl.New(nil).Require(ctx, permission.AccountSearch)).To(MatchError(authorization.ErrUnavailable))
	})
})

package role_test

import (
	"context"
	"testing"

	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	schoolgateway "grovehub/internal/session"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/constants/persona"
	roledto "grovehub/services/identity/internal/app/dto/role"
	"grovehub/services/identity/internal/app/entity"
	rolerepo "grovehub/services/identity/internal/app/repository/role"
	"grovehub/services/identity/internal/app/services/authorization"
	roleservice "grovehub/services/identity/internal/app/services/role"
	roleimpl "grovehub/services/identity/internal/app/services/role/impl"
	"grovehub/services/identity/internal/app/services/security/principal"
	rolevalidator "grovehub/services/identity/internal/app/validators/role"
)

const actorID = "11111111-1111-4111-8111-111111111111"
const roleID = "22222222-2222-4222-8222-222222222222"

func TestRoleService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Role Service Suite")
}

type authorizer struct {
	code permission.Code
	err  error
	grants []permission.Code
}

func (auth *authorizer) Require(ctx context.Context, code permission.Code) error {
	auth.code = code
	return auth.err
}

func (auth *authorizer) Effective(ctx context.Context) ([]permission.Code, error) {
	return auth.grants, auth.err
}

type repository struct {
	draft rolerepo.Draft
	update rolerepo.MetadataUpdate
	query coresearch.Query[string]
	permissions rolerepo.PermissionSet
	tenant string
	id string
	calls int
	err error
}

func (repo *repository) CreateDraft(ctx context.Context, draft rolerepo.Draft) (rolerepo.Detail, error) {
	repo.calls++
	repo.draft = draft
	return rolerepo.Detail{Role: entity.Role{RoleID: roleID, Name: draft.Name, Status: entity.RoleInactive, Version: 1}, Personas: draft.Personas}, repo.err
}

func (repo *repository) FindByID(ctx context.Context, tenantID string, roleID string) (rolerepo.Detail, error) {
	repo.calls++
	repo.tenant = tenantID
	repo.id = roleID
	return rolerepo.Detail{Role: entity.Role{RoleID: roleID, Status: entity.RoleInactive}}, repo.err
}

func (repo *repository) UpdateMetadata(ctx context.Context, update rolerepo.MetadataUpdate) (rolerepo.Detail, error) {
	repo.calls++
	repo.update = update
	return rolerepo.Detail{Role: entity.Role{RoleID: update.RoleID, Name: update.Name, Version: update.Version + 1}, Personas: update.Personas}, repo.err
}

func (repo *repository) Search(ctx context.Context, tenantID string, query coresearch.Query[string]) (coresearch.QueryResult[rolerepo.Detail], error) {
	repo.calls++
	repo.tenant = tenantID
	repo.query = query
	total := int64(1)
	return coresearch.QueryResult[rolerepo.Detail]{Items: []rolerepo.Detail{{Role: entity.Role{RoleID: roleID, Name: "Teacher"}}}, TotalRecords: &total}, repo.err
}

func (repo *repository) Permissions(ctx context.Context, tenantID string, roleID string) (rolerepo.PermissionSet, error) {
	repo.calls++
	repo.tenant = tenantID
	repo.id = roleID
	return repo.permissions, repo.err
}

var repo *repository
var auth *authorizer
var service roleservice.Service
var ctx context.Context

var _ = BeforeEach(func() {
	repo = &repository{}
	auth = &authorizer{}
	service = roleimpl.New(repo, auth)
	ctx = schoolgateway.WithPrincipal(context.Background(), schoolgateway.Principal{AccountID: actorID, TenantID: "school-a", SessionID: "session", Admin: true})
})

var _ = Describe("CreateDraft", func() {
	It("requires role create and derives tenant and actor from principal", func() {
		result, err := service.CreateDraft(ctx, roledto.CreateDraft{Name: " Teacher ", Description: " Test ", Personas: []persona.Code{persona.Teacher}})
		Expect(err).NotTo(HaveOccurred())
		Expect(auth.code).To(Equal(permission.RoleCreate))
		Expect(repo.draft).To(Equal(rolerepo.Draft{TenantID: "school-a", ActorID: actorID, Name: "Teacher", Description: "Test", Personas: []persona.Code{persona.Teacher}}))
		Expect(result.Status).To(Equal("INACTIVE"))
		Expect(result.Version).To(Equal(int64(1)))
	})
	It("rejects unauthenticated, denied and invalid requests without writes", func() {
		request := roledto.CreateDraft{Name: "Teacher", Personas: []persona.Code{persona.Teacher}}
		_, err := service.CreateDraft(context.Background(), request)
		Expect(err).To(MatchError(principal.ErrAuthenticationRequired))
		auth.err = authorization.ErrForbidden
		_, err = service.CreateDraft(ctx, request)
		Expect(err).To(MatchError(authorization.ErrForbidden))
		auth.err = nil
		_, err = service.CreateDraft(ctx, roledto.CreateDraft{})
		Expect(err).To(MatchError(rolevalidator.ErrInvalidDraft))
		Expect(repo.calls).To(BeZero())
	})
	It("does not return a result on a persistence failure", func() {
		repo.err = rolerepo.ErrConflict
		result, err := service.CreateDraft(ctx, roledto.CreateDraft{Name: "Teacher", Personas: []persona.Code{persona.Teacher}})
		Expect(err).To(MatchError(rolerepo.ErrConflict))
		Expect(result).To(Equal(roledto.Detail{}))
	})
})

var _ = Describe("Get", func() {
	It("requires role view and reads only the principal tenant", func() {
		result, err := service.Get(ctx, roleID)
		Expect(err).NotTo(HaveOccurred())
		Expect(auth.code).To(Equal(permission.RoleView))
		Expect(repo.tenant).To(Equal("school-a"))
		Expect(repo.id).To(Equal(roleID))
		Expect(result.RoleID).To(Equal(roleID))
		Expect(result.Personas).To(BeEmpty())
	})
	It("fails closed when dependencies are absent or authority is denied", func() {
		_, err := roleimpl.New(nil, auth).Get(ctx, roleID)
		Expect(err).To(MatchError(authorization.ErrUnavailable))
		_, err = roleimpl.New(repo, nil).Get(ctx, roleID)
		Expect(err).To(MatchError(authorization.ErrUnavailable))
		auth.err = authorization.ErrForbidden
		_, err = service.Get(ctx, roleID)
		Expect(err).To(MatchError(authorization.ErrForbidden))
		auth.err = nil
		_, err = service.Get(ctx, "bad")
		Expect(err).To(MatchError(rolevalidator.ErrInvalidID))
		Expect(repo.calls).To(BeZero())
	})
	It("preserves a tenant-scoped not found error", func() {
		repo.err = rolerepo.ErrNotFound
		result, err := service.Get(ctx, roleID)
		Expect(err).To(MatchError(rolerepo.ErrNotFound))
		Expect(result).To(Equal(roledto.Detail{}))
	})
})

var _ = Describe("UpdateMetadata", func() {
	It("requires role update and passes scoped metadata and expected version", func() {
		result, err := service.UpdateMetadata(ctx, roleID, roledto.UpdateMetadata{Name: " Teacher ", Personas: []persona.Code{persona.Teacher}, Version: 4})
		Expect(err).NotTo(HaveOccurred())
		Expect(auth.code).To(Equal(permission.RoleUpdate))
		Expect(repo.update.TenantID).To(Equal("school-a"))
		Expect(repo.update.ActorID).To(Equal(actorID))
		Expect(repo.update.Name).To(Equal("Teacher"))
		Expect(repo.update.Version).To(Equal(int64(4)))
		Expect(result.Version).To(Equal(int64(5)))
	})
	It("rejects forbidden and invalid updates without persistence", func() {
		request := roledto.UpdateMetadata{Name: "Teacher", Personas: []persona.Code{persona.Teacher}, Version: 1}
		auth.err = authorization.ErrForbidden
		_, err := service.UpdateMetadata(ctx, roleID, request)
		Expect(err).To(MatchError(authorization.ErrForbidden))
		auth.err = nil
		request.Version = 0
		_, err = service.UpdateMetadata(ctx, roleID, request)
		Expect(err).To(MatchError(rolevalidator.ErrInvalidVersion))
		Expect(repo.calls).To(BeZero())
	})
	It("returns no success response on a stale write", func() {
		repo.err = rolerepo.ErrVersionConflict
		result, err := service.UpdateMetadata(ctx, roleID, roledto.UpdateMetadata{Name: "Teacher", Personas: []persona.Code{persona.Teacher}, Version: 1})
		Expect(err).To(MatchError(rolerepo.ErrVersionConflict))
		Expect(result).To(Equal(roledto.Detail{}))
	})
})

var _ = Describe("Search", func() {
	It("uses verified tenant and role view permission", func() {
		result, err := service.Search(ctx, coresearch.Request{PageInfo: coresearch.PageRequest{PageNumber: 1}})
		Expect(err).NotTo(HaveOccurred())
		Expect(auth.code).To(Equal(permission.RoleView))
		Expect(repo.tenant).To(Equal("school-a"))
		Expect(repo.query.Limit).To(Equal(25))
		Expect(result.Items).To(HaveLen(1))
		Expect(*result.PageInfo.TotalRecords).To(Equal(int64(1)))
	})
	It("rejects denied or unknown filters before persistence", func() {
		auth.err = authorization.ErrForbidden
		_, err := service.Search(ctx, coresearch.Request{})
		Expect(err).To(MatchError(authorization.ErrForbidden))
		auth.err = nil
		_, err = service.Search(ctx, coresearch.Request{Filters: []coresearch.Filter{{Field: "tenantId", Operator: coresearch.OperatorEquals, Values: []string{"school-b"}}}})
		Expect(err).To(MatchError(rolevalidator.ErrInvalidSearch))
		Expect(repo.calls).To(BeZero())
	})
})

var _ = Describe("Permissions", func() {
	It("requires mapping view and scopes the read to the actor tenant", func() {
		repo.permissions = rolerepo.PermissionSet{RoleID: roleID, Version: 3}
		result, err := service.Permissions(ctx, roleID)
		Expect(err).NotTo(HaveOccurred())
		Expect(auth.code).To(Equal(permission.RolePermissionView))
		Expect(repo.tenant).To(Equal("school-a"))
		Expect(result.Codes).To(BeEmpty())
		Expect(result.Version).To(Equal(int64(3)))
	})
})

var _ = Describe("PreviewPermissions", func() {
	It("denies a legacy admin without effective mapping management authority", func() {
		_, err := service.PreviewPermissions(ctx, roleID, roledto.PreviewPermissions{Version: 1, Codes: []permission.Code{}})
		Expect(err).To(MatchError(authorization.ErrForbidden))
		Expect(repo.calls).To(BeZero())
	})
	It("validates the version and marks changes as requiring approval", func() {
		auth.grants = []permission.Code{permission.RolePermissionManage, permission.AccountView}
		repo.permissions = rolerepo.PermissionSet{RoleID: roleID, Version: 2}
		request := roledto.PreviewPermissions{Version: 1, Codes: []permission.Code{permission.AccountView}}
		_, err := service.PreviewPermissions(ctx, roleID, request)
		Expect(err).To(MatchError(rolerepo.ErrVersionConflict))
		request.Version = 2
		result, err := service.PreviewPermissions(ctx, roleID, request)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.Added).To(Equal(request.Codes))
		Expect(result.ApprovalRequired).To(BeTrue())
		Expect(repo.tenant).To(Equal("school-a"))
	})
})

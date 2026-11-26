package registry_test

import (
	"context"
	"errors"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/services/identity/internal/app/constants/feature"
	"grovehub/services/identity/internal/app/constants/permission"
	registryrepo "grovehub/services/identity/internal/app/repository/registry"
	"grovehub/services/identity/internal/app/services/registry"
	"grovehub/services/identity/internal/app/services/registry/helpers"
	registryimpl "grovehub/services/identity/internal/app/services/registry/impl"
)

func TestRegistryService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Registry Service Suite")
}

type recorder struct {
	calls int
	batch registryrepo.Batch
	err   error
}

func (repo *recorder) Register(ctx context.Context, batch registryrepo.Batch) (registryrepo.Result, error) {
	repo.calls++
	repo.batch = batch
	if repo.err != nil {
		return registryrepo.Result{}, repo.err
	}
	return registryrepo.Result{FeaturesCreated: len(batch.Features), PermissionsCreated: len(batch.Permissions)}, nil
}

var _ = Describe("PrepareBatch", func() {
	It("resolves canonical metadata and deterministic order without modifying input", func() {
		request := registry.Request{
			Features:    []feature.Code{feature.RoleManagement, feature.AccountManagement},
			Permissions: []permission.Code{permission.RoleView, permission.AccountSearch, permission.AccountSelfView},
		}
		batch, err := helpers.PrepareBatch(request)
		Expect(err).NotTo(HaveOccurred())
		Expect(batch.Features[0].Code).To(Equal(feature.AccountManagement))
		Expect(batch.Permissions[0].Code).To(Equal(permission.AccountSearch))
		Expect(batch.Permissions[0].Feature).To(Equal(feature.AccountManagement))
		Expect(batch.Permissions[1].Feature).To(BeEmpty())
		Expect(batch.Permissions[1].Name).To(Equal("View own account"))
		Expect(request.Features[0]).To(Equal(feature.RoleManagement))
		Expect(request.Permissions[0]).To(Equal(permission.RoleView))
	})

	It("rejects empty, oversized, unknown and duplicate input", func() {
		cases := []registry.Request{
			{},
			{Features: make([]feature.Code, 101)},
			{Features: []feature.Code{"UNKNOWN"}},
			{Permissions: []permission.Code{"UNKNOWN"}},
			{Features: []feature.Code{feature.RoleManagement, feature.RoleManagement}},
			{Permissions: []permission.Code{permission.AccountSearch, permission.AccountSearch}},
		}
		for _, request := range cases {
			batch, err := helpers.PrepareBatch(request)
			Expect(err).To(MatchError(registry.ErrInvalidBatch))
			Expect(batch).To(Equal(registryrepo.Batch{}))
		}
	})

	It("keeps every known permission binding within the code-owned feature vocabulary", func() {
		seen := make(map[permission.Code]bool)
		for _, definition := range permission.Catalog() {
			Expect(seen[definition.Code]).To(BeFalse())
			seen[definition.Code] = true
			Expect(definition.Name).NotTo(BeEmpty())
			if definition.Feature != "" {
				Expect(feature.Known(definition.Feature)).To(BeTrue())
			}
		}
		Expect(seen).To(HaveLen(22))
	})
})

var _ = Describe("Register", func() {
	It("rejects invalid input before persistence and forwards only canonical definitions", func() {
		repo := &recorder{}
		service := registryimpl.New(repo)
		_, err := service.Register(context.Background(), registry.Request{Permissions: []permission.Code{"UNKNOWN"}})
		Expect(err).To(MatchError(registry.ErrInvalidBatch))
		Expect(repo.calls).To(BeZero())
		result, err := service.Register(context.Background(), registry.Request{Features: []feature.Code{feature.RoleManagement}, Permissions: []permission.Code{permission.RoleView}})
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(registryrepo.Result{FeaturesCreated: 1, PermissionsCreated: 1}))
		Expect(repo.calls).To(Equal(1))
		Expect(repo.batch.Permissions[0].Feature).To(Equal(feature.RoleManagement))
		repo.err = errors.New("transaction failed")
		result, err = service.Register(context.Background(), registry.Request{Permissions: []permission.Code{permission.RoleView}})
		Expect(err).To(MatchError(repo.err))
		Expect(result).To(Equal(registryrepo.Result{}))
	})
})

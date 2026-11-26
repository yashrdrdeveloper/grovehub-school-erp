package delegation_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/authorization/operations/delegation"
)

func TestDelegation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Delegation Suite")
}

var _ = Describe("Preview", func() {
	It("returns deterministic changes without mutating inputs", func() {
		held := []permission.Code{permission.RolePermissionManage, permission.AccountView, permission.AccountSearch, permission.AccountUpdate}
		current := []permission.Code{permission.AccountUpdate}
		requested := []permission.Code{permission.AccountView, permission.AccountSearch}
		change, err := delegation.Preview(held, current, requested)
		Expect(err).NotTo(HaveOccurred())
		Expect(change.Added).To(Equal([]permission.Code{permission.AccountSearch, permission.AccountView}))
		Expect(change.Removed).To(Equal(current))
		Expect(requested[0]).To(Equal(permission.AccountView))
	})
	It("denies missing manager capability, unheld grants, unheld removals and restricted permissions", func() {
		_, err := delegation.Preview(nil, nil, nil)
		Expect(err).To(MatchError(authorization.ErrForbidden))
		held := []permission.Code{permission.RolePermissionManage}
		_, err = delegation.Preview(held, nil, []permission.Code{permission.AccountView})
		Expect(err).To(MatchError(authorization.ErrForbidden))
		_, err = delegation.Preview(held, []permission.Code{permission.AccountView}, nil)
		Expect(err).To(MatchError(authorization.ErrForbidden))
		for _, code := range []permission.Code{permission.PermissionRegister, permission.FeatureRegister, permission.FeatureManage} {
			_, err = delegation.Preview(append(held, code), nil, []permission.Code{code})
			Expect(err).To(MatchError(authorization.ErrForbidden))
		}
	})
	It("rejects unknown or duplicate codes and supports empty replacement", func() {
		held := []permission.Code{permission.RolePermissionManage, permission.AccountView}
		_, err := delegation.Preview(held, nil, []permission.Code{"UNKNOWN"})
		Expect(err).To(MatchError(delegation.ErrInvalidSet))
		_, err = delegation.Preview(held, nil, []permission.Code{permission.AccountView, permission.AccountView})
		Expect(err).To(MatchError(delegation.ErrInvalidSet))
		change, err := delegation.Preview(held, []permission.Code{permission.AccountView}, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(change.Removed).To(Equal([]permission.Code{permission.AccountView}))
		change, err = delegation.Preview(held, nil, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(change.Added).To(BeEmpty())
		Expect(change.Removed).To(BeEmpty())
	})
})

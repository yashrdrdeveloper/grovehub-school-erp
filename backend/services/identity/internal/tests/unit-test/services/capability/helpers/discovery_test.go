package helpers_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/services/identity/internal/app/constants/feature"
	"grovehub/services/identity/internal/app/constants/permission"
	repo "grovehub/services/identity/internal/app/repository/capability"
	"grovehub/services/identity/internal/app/services/capability/helpers"
)

func TestCapabilityHelpers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Capability Helpers Suite")
}

var _ = Describe("Actions", func() {
	It("does not advertise blocked or absent routes even for an actor holding their permissions", func() {
		actions := helpers.Actions([]permission.Code{permission.AccountCreate, permission.RoleCreate, permission.RolePermissionManage})
		Expect(actions["accountCreate"]).To(BeFalse())
		Expect(actions["roleCreate"]).To(BeFalse())
		Expect(actions["rolePermissionsSubmit"]).To(BeTrue())
		Expect(actions["permissionCatalog"]).To(BeFalse())
	})
})

var _ = Describe("Catalog", func() {
	It("excludes reserved, unknown and wrongly bound permissions and only marks held grants delegable", func() {
		rows := []repo.Entry{
			{Code: permission.RoleView, Feature: feature.RoleManagement},
			{Code: permission.RoleCreate, Feature: feature.RoleManagement},
			{Code: permission.RoleUpdate, Feature: feature.AccountManagement},
			{Code: permission.FeatureManage},
			{Code: "UNKNOWN"},
		}
		result := helpers.Catalog(rows, []permission.Code{permission.RoleView, permission.RolePermissionManage})
		Expect(result.Permissions).To(HaveLen(2))
		Expect(result.Permissions[0].CanDelegate).To(BeTrue())
		Expect(result.Permissions[1].CanDelegate).To(BeFalse())
		Expect(helpers.Catalog(rows, []permission.Code{permission.RoleView}).Permissions[0].CanDelegate).To(BeFalse())
	})
})

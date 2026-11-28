package manifest_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/services/rbacmigration"
	"grovehub/services/identity/internal/app/services/rbacmigration/operations/manifest"
)

func TestManifest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bootstrap Manifest Suite")
}

var _ = Describe("Validate", func() {
	It("requires explicit recovery grants on roles held exclusively by the nominated administrator", func() {
		admin := "11111111-1111-4111-8111-111111111111"
		request := rbacmigration.Manifest{
			AdministratorID: admin,
			Plan: rbacmigration.Plan{Roles: []rbacmigration.RolePlan{{RoleID: "role"}}, Assignments: []rbacmigration.AssignmentPlan{{AccountID: admin, PrimaryRoleID: "role"}}},
			Grants: []rbacmigration.RoleGrant{{RoleID: "role", Codes: permission.AdministratorRecoveryCapabilities()}},
		}
		Expect(manifest.Validate(request)).To(Succeed())
		request.Grants[0].Codes = []permission.Code{permission.RolePermissionManage}
		Expect(manifest.Validate(request)).To(MatchError(rbacmigration.ErrAdministrator))
		request.Grants[0].Codes = append(permission.AdministratorRecoveryCapabilities(), permission.FeatureManage)
		Expect(manifest.Validate(request)).To(MatchError(rbacmigration.ErrInvalidManifest))
		request.Grants[0].Codes = permission.AdministratorRecoveryCapabilities()
		request.Plan.Assignments = append(request.Plan.Assignments, rbacmigration.AssignmentPlan{AccountID: "other", PrimaryRoleID: "role"})
		Expect(manifest.Validate(request)).To(MatchError(rbacmigration.ErrAdministrator))
	})
})

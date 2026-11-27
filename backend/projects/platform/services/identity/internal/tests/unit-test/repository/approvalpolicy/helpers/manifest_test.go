package helpers_test

import (
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	kind "grovehub/services/identity/internal/app/constants/approval"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/repository/approvalpolicy"
	"grovehub/services/identity/internal/app/repository/approvalpolicy/helpers"
)

func TestApprovalPolicyManifest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Approval Policy Manifest Suite")
}

var _ = Describe("Validate", func() {
	It("requires all four supported types with explicit ordered steps and versions", func() {
		const author = "11111111-1111-4111-8111-111111111111"
		digest := strings.Repeat("a",64)
		steps := []approvalpolicy.Step{{Permission: permission.AccountRoleAssign}}
		manifest := approvalpolicy.Manifest{TenantID: "school", AuthorID: author, Policies: []approvalpolicy.Revision{
			{Type: kind.AccountLifecycleChange, Steps: steps},
			{Type: kind.RoleLifecycleChange, Steps: steps},
			{Type: kind.RolePermissionChange, Steps: steps},
			{Type: kind.RoleAssignmentChange, Steps: steps},
		}}
		Expect(helpers.Validate(manifest, author, digest)).To(Succeed())
		manifest.Policies[3].Type = kind.RoleLifecycleChange
		Expect(helpers.Validate(manifest, author, digest)).To(MatchError(approvalpolicy.ErrInvalid))
		manifest.Policies[3].Type = kind.RoleAssignmentChange
		manifest.Policies[3].ExpectedVersion = -1
		Expect(helpers.Validate(manifest, author, digest)).To(MatchError(approvalpolicy.ErrInvalid))
		manifest.Policies[3].ExpectedVersion = 0
		manifest.Policies[3].Steps = []approvalpolicy.Step{{Permission: permission.FeatureManage}}
		Expect(helpers.Validate(manifest, author, digest)).To(MatchError(approvalpolicy.ErrInvalid))
	})
})

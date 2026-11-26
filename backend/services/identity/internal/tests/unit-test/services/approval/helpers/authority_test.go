package helpers_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	kind "grovehub/services/identity/internal/app/constants/approval"
	"grovehub/services/identity/internal/app/constants/permission"
	accountdto "grovehub/services/identity/internal/app/dto/accountrole"
	dto "grovehub/services/identity/internal/app/dto/approval"
	"grovehub/services/identity/internal/app/entity"
	repo "grovehub/services/identity/internal/app/repository/approval"
	"grovehub/services/identity/internal/app/services/approval/helpers"
)

var _ = Describe("NormalizePermissions", func() {
	It("sorts codes, preserves explicit empty sets and rejects duplicates or reserved grants", func() {
		intent := dto.PermissionIntent{RoleID: target, Version: 1, Codes: []permission.Code{permission.RoleView, permission.RoleCreate}}
		normalized, err := helpers.NormalizePermissions(intent)
		Expect(err).NotTo(HaveOccurred())
		Expect(normalized.Codes).To(Equal([]permission.Code{permission.RoleCreate, permission.RoleView}))
		intent.Codes = []permission.Code{}
		normalized, err = helpers.NormalizePermissions(intent)
		Expect(err).NotTo(HaveOccurred())
		Expect(normalized.Codes).NotTo(BeNil())
		intent.Codes = nil
		_, err = helpers.NormalizePermissions(intent)
		Expect(err).To(MatchError(repo.ErrInvalid))
		intent.Codes = []permission.Code{permission.RoleView, permission.RoleView}
		_, err = helpers.NormalizePermissions(intent)
		Expect(err).To(MatchError(repo.ErrInvalid))
		intent.Codes = []permission.Code{permission.FeatureManage}
		_, err = helpers.NormalizePermissions(intent)
		Expect(err).To(MatchError(repo.ErrInvalid))
	})
})

var _ = Describe("NormalizeAssignment", func() {
	It("requires a version stamp for all proposed roles and sorts role versions", func() {
		other := "22222222-2222-4222-8222-222222222222"
		intent := dto.AssignmentIntent{AccountID: target, Version: 2, PrimaryRoleID: other, SecondaryRoleIDs: []string{target}, RoleVersions: []accountdto.RoleVersion{{RoleID: other, Version: 2}, {RoleID: target, Version: 1}}}
		normalized, err := helpers.NormalizeAssignment(intent)
		Expect(err).NotTo(HaveOccurred())
		Expect(normalized.RoleVersions[0].RoleID).To(Equal(target))
		intent.RoleVersions = intent.RoleVersions[:1]
		_, err = helpers.NormalizeAssignment(intent)
		Expect(err).To(MatchError(repo.ErrInvalid))
	})
})

var _ = Describe("VerifyIntent", func() {
	It("accepts JSONB key reordering but rejects changed versions and unknown fields", func() {
		intent := dto.PermissionIntent{RoleID: target, Version: 1, Codes: []permission.Code{permission.RoleView}}
		_, digest, err := helpers.EncodeIntent(intent)
		Expect(err).NotTo(HaveOccurred())
		request := entity.ApprovalRequest{RequestType: kind.RolePermissionChange, PayloadSHA256: digest, Payload: json.RawMessage(`{"version":1,"permissionCodes":["IDENTITY_ROLE_VIEW"],"roleId":"11111111-1111-4111-8111-111111111111"}`)}
		Expect(helpers.VerifyIntent(request)).To(Succeed())
		request.Payload = json.RawMessage(`{"version":2,"permissionCodes":["IDENTITY_ROLE_VIEW"],"roleId":"11111111-1111-4111-8111-111111111111"}`)
		Expect(helpers.VerifyIntent(request)).To(MatchError(repo.ErrConflict))
		request.Payload = json.RawMessage(`{"version":1,"permissionCodes":["IDENTITY_ROLE_VIEW"],"roleId":"11111111-1111-4111-8111-111111111111","extra":true}`)
		Expect(helpers.VerifyIntent(request)).To(MatchError(repo.ErrInvalid))
	})
})

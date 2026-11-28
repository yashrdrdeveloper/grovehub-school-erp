package preview_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/services/identity/internal/app/constants/permission"
	dto "grovehub/services/identity/internal/app/dto/accountrole"
	"grovehub/services/identity/internal/app/entity"
	repo "grovehub/services/identity/internal/app/repository/accountrole"
	"grovehub/services/identity/internal/app/services/accountrole/operations/preview"
	"grovehub/services/identity/internal/app/services/authorization"
)

func TestAssignmentPreview(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Account Role Assignment Preview Suite")
}

var _ = Describe("RoleIDs", func() {
	It("deduplicates and sorts old and new roles", func() {
		Expect(preview.RoleIDs(repo.Assignment{PrimaryRoleID: "b", SecondaryRoleIDs: []string{"c"}}, dto.Replace{PrimaryRoleID: "a", SecondaryRoleIDs: []string{"b"}})).To(Equal([]string{"a", "b", "c"}))
	})
})

var _ = Describe("Evaluate", func() {
	var current repo.Assignment
	var request dto.Replace
	var candidates []repo.Candidate
	var effective []permission.Code

	BeforeEach(func() {
		current = repo.Assignment{AccountID: "account", Status: entity.Active, Version: 2, PrimaryRoleID: "old"}
		request = dto.Replace{Version: 2, PrimaryRoleID: "new", SecondaryRoleIDs: []string{}}
		candidates = []repo.Candidate{
			{RoleID: "old", Status: entity.RoleActive, Version: 1, Codes: []permission.Code{permission.AccountSearch}},
			{RoleID: "new", Status: entity.RoleActive, Version: 3, Codes: []permission.Code{permission.AccountView}},
		}
		effective = []permission.Code{permission.AccountRoleAssign, permission.AccountSearch, permission.AccountView}
	})
	It("returns version-bound approval intent for eligible roles", func() {
		result, err := preview.Evaluate(current, request, candidates, effective)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.ApprovalRequired).To(BeTrue())
		Expect(result.RoleVersions).To(Equal([]dto.RoleVersion{{RoleID: "new", Version: 3}, {RoleID: "old", Version: 1}}))
	})
	It("rejects missing authority over current or proposed grants", func() {
		_, err := preview.Evaluate(current, request, candidates, effective[:2])
		Expect(err).To(MatchError(authorization.ErrForbidden))
		_, err = preview.Evaluate(current, request, candidates, []permission.Code{permission.AccountRoleAssign, permission.AccountView})
		Expect(err).To(MatchError(authorization.ErrForbidden))
	})
	It("rejects inactive targets, stale versions, and suspended accounts", func() {
		candidates[1].Status = entity.RoleInactive
		_, err := preview.Evaluate(current, request, candidates, effective)
		Expect(err).To(MatchError(repo.ErrIneligible))
		request.Version = 1
		_, err = preview.Evaluate(current, request, candidates, effective)
		Expect(err).To(MatchError(repo.ErrVersionConflict))
		request.Version = 2
		current.Status = entity.Suspended
		_, err = preview.Evaluate(current, request, candidates, effective)
		Expect(err).To(MatchError(repo.ErrIneligible))
	})
})

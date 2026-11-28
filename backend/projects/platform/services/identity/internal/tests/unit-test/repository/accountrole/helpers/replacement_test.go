package helpers_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/repository/accountrole"
	"grovehub/services/identity/internal/app/repository/accountrole/helpers"
)

func TestReplacement(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Account Role Replacement Suite")
}

const first = "11111111-1111-4111-8111-111111111111"
const second = "22222222-2222-4222-8222-222222222222"

var _ = Describe("NormalizeReplacement", func() {
	It("rejects duplicate IDs and invalid version stamps", func() {
		request := accountrole.Replacement{TenantID: "school", ActorID: first, AccountID: first, Version: 1, PrimaryRoleID: first, SecondaryRoleIDs: []string{first}}
		_, err := helpers.NormalizeReplacement(request)
		Expect(err).To(MatchError(accountrole.ErrInvalidAssignment))
		request.SecondaryRoleIDs = []string{second}
		request.RoleVersions = []accountrole.RoleVersion{{RoleID: first, Version: 0}}
		_, err = helpers.NormalizeReplacement(request)
		Expect(err).To(MatchError(accountrole.ErrInvalidAssignment))
		request.RoleVersions[0].Version = 1
		actual, err := helpers.NormalizeReplacement(request)
		Expect(err).NotTo(HaveOccurred())
		Expect(actual).To(Equal(request))
	})
})

var _ = Describe("ValidateCandidates", func() {
	It("rejects stale roles and ineligible targets", func() {
		current := accountrole.Assignment{Version: 1, Status: entity.Active, PrimaryRoleID: first}
		request := accountrole.Replacement{Version: 1, PrimaryRoleID: second, RoleVersions: []accountrole.RoleVersion{{RoleID: first, Version: 1}, {RoleID: second, Version: 2}}}
		candidates := []accountrole.Candidate{{RoleID: first, Version: 1}, {RoleID: second, Version: 2, Status: entity.RoleActive, Codes: []permission.Code{permission.AccountView}}}
		Expect(helpers.ValidateCandidates(current, request, candidates)).To(Succeed())
		candidates[1].Version = 3
		Expect(helpers.ValidateCandidates(current, request, candidates)).To(MatchError(accountrole.ErrVersionConflict))
		candidates[1].Version = 2
		candidates[1].Status = entity.RoleInactive
		Expect(helpers.ValidateCandidates(current, request, candidates)).To(MatchError(accountrole.ErrIneligible))
	})
})

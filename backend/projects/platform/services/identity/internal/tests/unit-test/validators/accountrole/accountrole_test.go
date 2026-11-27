package accountrole_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	dto "grovehub/services/identity/internal/app/dto/accountrole"
	repo "grovehub/services/identity/internal/app/repository/accountrole"
	validator "grovehub/services/identity/internal/app/validators/accountrole"
)

func TestAccountRoleValidator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Account Role Validator Suite")
}

const primary = "11111111-1111-4111-8111-111111111111"
const secondary = "22222222-2222-4222-8222-222222222222"

var _ = Describe("NormalizeReplace", func() {
	It("requires one primary and rejects duplicate assignments", func() {
		_, err := validator.NormalizeReplace(dto.Replace{Version: 1, PrimaryRoleID: primary, SecondaryRoleIDs: []string{primary}})
		Expect(err).To(MatchError(repo.ErrInvalidAssignment))
		_, err = validator.NormalizeReplace(dto.Replace{Version: 1, PrimaryRoleID: primary, SecondaryRoleIDs: []string{secondary, secondary}})
		Expect(err).To(MatchError(repo.ErrInvalidAssignment))
		_, err = validator.NormalizeReplace(dto.Replace{Version: 1})
		Expect(err).To(MatchError(repo.ErrInvalidAssignment))
	})
	It("returns normalized UUIDs without changing the caller input", func() {
		request := dto.Replace{Version: 2, PrimaryRoleID: primary, SecondaryRoleIDs: []string{secondary}}
		result, err := validator.NormalizeReplace(request)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(request))
		result.SecondaryRoleIDs[0] = primary
		Expect(request.SecondaryRoleIDs[0]).To(Equal(secondary))
	})
})

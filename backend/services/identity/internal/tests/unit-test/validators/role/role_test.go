package role_test

import (
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/services/identity/internal/app/constants/persona"
	roledto "grovehub/services/identity/internal/app/dto/role"
	rolevalidator "grovehub/services/identity/internal/app/validators/role"
)

func TestRoleValidator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Role Validator Suite")
}

var _ = Describe("NormalizeCreateDraft", func() {
	It("normalizes whitespace while preserving input", func() {
		input := roledto.CreateDraft{Name: " Teacher ", Description: " Example ", Personas: []persona.Code{persona.Teacher}}
		actual, err := rolevalidator.NormalizeCreateDraft(input)
		Expect(err).NotTo(HaveOccurred())
		Expect(actual.Name).To(Equal("Teacher"))
		Expect(actual.Description).To(Equal("Example"))
		Expect(input.Name).To(Equal(" Teacher "))
	})
	It("rejects invalid names, descriptions and personas", func() {
		for _, input := range []roledto.CreateDraft{
			{},
			{Name: "Teacher"},
			{Name: strings.Repeat("a", 121), Personas: []persona.Code{persona.Teacher}},
			{Name: "Teacher", Description: strings.Repeat("a", 1001), Personas: []persona.Code{persona.Teacher}},
			{Name: "Teacher", Personas: []persona.Code{"UNKNOWN"}},
			{Name: "Teacher", Personas: []persona.Code{persona.Teacher, persona.Teacher}},
			{Name: "Teacher", Description: "\x00", Personas: []persona.Code{persona.Teacher}},
			{Name: string([]byte{0xff}), Personas: []persona.Code{persona.Teacher}},
		} {
			_, err := rolevalidator.NormalizeCreateDraft(input)
			Expect(err).To(MatchError(rolevalidator.ErrInvalidDraft))
		}
	})
})

var _ = Describe("NormalizeID", func() {
	It("normalizes valid UUIDs and rejects malformed or nil UUIDs", func() {
		actual, err := rolevalidator.NormalizeID("AAAAAAAA-AAAA-4AAA-8AAA-AAAAAAAAAAAA")
		Expect(err).NotTo(HaveOccurred())
		Expect(actual).To(Equal("aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"))
		for _, value := range []string{"", "bad", "00000000-0000-0000-0000-000000000000"} {
			_, err := rolevalidator.NormalizeID(value)
			Expect(err).To(MatchError(rolevalidator.ErrInvalidID))
		}
	})
})

var _ = Describe("NormalizeUpdateMetadata", func() {
	It("normalizes complete metadata and preserves the expected version", func() {
		result, err := rolevalidator.NormalizeUpdateMetadata(roledto.UpdateMetadata{Name: " Teacher ", Personas: []persona.Code{persona.Teacher}, Version: 2})
		Expect(err).NotTo(HaveOccurred())
		Expect(result.Name).To(Equal("Teacher"))
		Expect(result.Version).To(Equal(int64(2)))
	})
	It("rejects nonpositive versions and invalid metadata", func() {
		_, err := rolevalidator.NormalizeUpdateMetadata(roledto.UpdateMetadata{})
		Expect(err).To(MatchError(rolevalidator.ErrInvalidVersion))
		_, err = rolevalidator.NormalizeUpdateMetadata(roledto.UpdateMetadata{Version: 1})
		Expect(err).To(MatchError(rolevalidator.ErrInvalidDraft))
	})
})

package entity_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/services/foundation/entity"
)

var _ = Describe("GID", func() {
	It("parses a canonical UUID", func() {
		gid, err := entity.ParseGID(validGID)

		Expect(err).NotTo(HaveOccurred())
		Expect(gid.String()).To(Equal(validGID))
		Expect(gid.IsZero()).To(BeFalse())
	})

	It("normalizes UUID casing", func() {
		gid, err := entity.ParseGID("01890F57-7B5C-7CC6-98B2-1F3D4E5A6B7C")

		Expect(err).NotTo(HaveOccurred())
		Expect(gid.String()).To(Equal(validGID))
	})

	It("rejects malformed UUID values", func() {
		_, err := entity.ParseGID("not-a-uuid")

		Expect(err).To(MatchError(entity.ErrInvalidGID))
	})

	It("rejects the zero UUID", func() {
		_, err := entity.ParseGID("00000000-0000-0000-0000-000000000000")

		Expect(err).To(MatchError(entity.ErrInvalidGID))
	})
})

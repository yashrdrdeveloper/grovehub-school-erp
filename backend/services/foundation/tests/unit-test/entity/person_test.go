package entity_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/services/foundation/entity"
)

var _ = Describe("Person", func() {
	It("creates an active person at version one", func() {
		gid, err := entity.ParseGID(validGID)
		Expect(err).NotTo(HaveOccurred())

		person, err := entity.NewPerson(gid, fixedTime)

		Expect(err).NotTo(HaveOccurred())
		Expect(person.GID).To(Equal(gid))
		Expect(person.Status).To(Equal(entity.PersonStatusActive))
		Expect(person.Version).To(Equal(int64(1)))
		Expect(person.CreatedAt).To(Equal(fixedTime))
		Expect(person.UpdatedAt).To(Equal(fixedTime))
		Expect(person.Validate()).To(Succeed())
	})

	It("keeps person status values owned by the person entity", func() {
		Expect(entity.PersonStatusActive.IsValid()).To(BeTrue())
		Expect(entity.PersonStatusDeactivated.IsValid()).To(BeTrue())
		Expect(entity.PersonStatus("UNKNOWN").IsValid()).To(BeFalse())
	})

	It("rejects a person without a gid", func() {
		_, err := entity.NewPerson(entity.GID{}, fixedTime)

		Expect(err).To(MatchError(entity.ErrPersonGIDRequired))
	})
})

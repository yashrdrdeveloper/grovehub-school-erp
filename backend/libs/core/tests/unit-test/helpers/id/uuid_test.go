package id_test

import (
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coreid "grovehub/libs/core/helpers/id"
)

var _ = Describe("UUIDv7 helper", func() {
	It("generates a UUID version 7 value", func() {
		value, err := coreid.NewUUIDv7()
		Expect(err).NotTo(HaveOccurred())

		parsed, err := uuid.Parse(value)
		Expect(err).NotTo(HaveOccurred())
		Expect(parsed.Version()).To(Equal(uuid.Version(7)))
	})
})

package uuid_test

import (
	"errors"
	"testing/iotest"

	googleuuid "github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/id/internal/src/operations/uuid"
)

var _ = Describe("NewUUIDv7", func() {
	It("generates distinct canonical UUIDv7 identifiers", func() {
		first, err := uuid.NewUUIDv7()
		Expect(err).NotTo(HaveOccurred())
		second, err := uuid.NewUUIDv7()
		Expect(err).NotTo(HaveOccurred())
		Expect(first).NotTo(Equal(second))
		for _, value := range []string{first, second} {
			parsed, err := googleuuid.Parse(value)
			Expect(err).NotTo(HaveOccurred())
			Expect(parsed.Version()).To(Equal(googleuuid.Version(7)))
			Expect(parsed.Variant()).To(Equal(googleuuid.RFC4122))
			Expect(parsed.String()).To(Equal(value))
		}
	})

	It("returns the original entropy failure without producing an identifier", func() {
		// The UUID entropy source is process-global. This suite runs serially
		// in an isolated test binary and restores the default after the spec.
		failure := errors.New("entropy unavailable")
		googleuuid.SetRand(iotest.ErrReader(failure))
		DeferCleanup(func() { googleuuid.SetRand(nil) })

		value, err := uuid.NewUUIDv7()
		Expect(err).To(BeIdenticalTo(failure))
		Expect(value).To(BeEmpty())
	})
})

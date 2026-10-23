package id

import (
	"bytes"
	"io"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UUIDv7", func() {
	It("generates a deterministic RFC 9562 UUIDv7 from controlled inputs", func() {
		random := bytes.NewReader([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})

		value, err := newUUIDv7(fixedUUIDv7Time, random)

		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal(expectedUUIDv7))
	})

	It("rejects timestamps before the Unix epoch", func() {
		_, err := newUUIDv7(time.UnixMilli(-1), bytes.NewReader(make([]byte, 10)))

		Expect(err).To(MatchError(ErrUUIDv7TimestampOutOfRange))
	})

	It("returns an error when secure randomness cannot be read", func() {
		_, err := newUUIDv7(fixedUUIDv7Time, bytes.NewReader([]byte{1}))

		Expect(err).To(MatchError(ContainSubstring(io.ErrUnexpectedEOF.Error())))
	})
})

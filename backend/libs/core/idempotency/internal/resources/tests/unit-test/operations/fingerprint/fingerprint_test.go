package fingerprint_test

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/idempotency/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/idempotency/internal/src/operations/fingerprint"
)

type invalidJSON struct{}

func (invalidJSON) MarshalJSON() ([]byte, error) {
	return nil, errors.New("cannot marshal")
}

var _ = Describe("Fingerprint", func() {
	It("preserves SHA-256 of the original JSON encoding", func() {
		value := struct {
			Name string `json:"name"`
		}{Name: "Ada"}
		digest := sha256.Sum256([]byte(`{"name":"Ada"}`))
		result, err := fingerprint.Fingerprint(value)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(hex.EncodeToString(digest[:])))
		Expect(result).To(HaveLen(64))
	})

	It("is deterministic for equivalent JSON maps and distinguishes changed values", func() {
		first, err := fingerprint.Fingerprint(map[string]int{"b": 2, "a": 1})
		Expect(err).NotTo(HaveOccurred())
		second, err := fingerprint.Fingerprint(map[string]int{"a": 1, "b": 2})
		Expect(err).NotTo(HaveOccurred())
		Expect(first).To(Equal(second))
		changed, err := fingerprint.Fingerprint(map[string]int{"a": 2, "b": 2})
		Expect(err).NotTo(HaveOccurred())
		Expect(changed).NotTo(Equal(first))
		nilDigest, err := fingerprint.Fingerprint(nil)
		Expect(err).NotTo(HaveOccurred())
		digest := sha256.Sum256([]byte("null"))
		Expect(nilDigest).To(Equal(hex.EncodeToString(digest[:])))
	})

	It("wraps JSON encoding failures with the fingerprint sentinel", func() {
		for _, value := range []any{make(chan int), math.NaN(), invalidJSON{}} {
			result, err := fingerprint.Fingerprint(value)
			Expect(result).To(BeEmpty())
			Expect(errors.Is(err, contracts.ErrFingerprint)).To(BeTrue())
		}
	})
})

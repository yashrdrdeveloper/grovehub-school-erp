package key_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/idempotency/internal/src/contracts"
	"grovehub/libs/core/idempotency/internal/src/operations/key"
)

var _ = Describe("ValidateKey", func() {
	It("trims only surrounding whitespace and accepts ASCII letters, digits and permitted punctuation", func() {
		for _, value := range []string{"ABCxyz019-_.:", "register:019-test", strings.Repeat("a", 128)} {
			result, err := key.ValidateKey("  " + value + " \t")
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(value))
		}
	})

	It("rejects empty, oversized or invalid keys without returning partial values", func() {
		for _, value := range []string{"", " \t", strings.Repeat("a", 129), "a b", "a\tb", "a\nb", "a/b", "a,b", "a@b", "é"} {
			result, err := key.ValidateKey(value)
			Expect(result).To(BeEmpty())
			Expect(err).To(BeIdenticalTo(contracts.ErrInvalidKey))
		}
	})
})

package helpers_test

import (
	"grovehub/libs/core/headercontext/internal/src/helpers"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ValidID", func() {
	It("accepts optional values, supported ASCII and the exact length limit", func() {
		for _, value := range []string{"", "aZ09-_.:/", strings.Repeat("a", 128)} {
			Expect(helpers.ValidID(value)).To(BeTrue())
		}
	})

	It("rejects oversized values, whitespace, non-ASCII and delimiters", func() {
		for _, value := range []string{strings.Repeat("a", 129), "a b", "a\n", "a\t", "é", "a,b", "a@b", "a\\b"} {
			Expect(helpers.ValidID(value)).To(BeFalse())
		}
	})
})

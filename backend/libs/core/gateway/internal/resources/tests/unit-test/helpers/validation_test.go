package helpers_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/gateway/internal/src/helpers"
)

var _ = Describe("ValidKeyID", func() {
	It("accepts existing valid key identifiers including the maximum length", func() {
		Expect(helpers.ValidKeyID("school-1")).To(BeTrue())
		Expect(helpers.ValidKeyID(strings.Repeat("a", 128))).To(BeTrue())
	})
	It("rejects empty, oversized or delimiter-containing identifiers", func() {
		for _, value := range []string{"", strings.Repeat("a", 129), "a.b", "a b", "a\tb", "a\rb", "a\nb", "a,b"} {
			Expect(helpers.ValidKeyID(value)).To(BeFalse())
		}
	})
})

var _ = Describe("ValidActorPair", func() {
	It("accepts both present or both absent and rejects partial delegation", func() {
		Expect(helpers.ValidActorPair("", "")).To(BeTrue())
		Expect(helpers.ValidActorPair("tenant", "actor")).To(BeTrue())
		Expect(helpers.ValidActorPair("tenant", "")).To(BeFalse())
		Expect(helpers.ValidActorPair("", "actor")).To(BeFalse())
	})
})

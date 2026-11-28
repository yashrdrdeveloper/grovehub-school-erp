package contracts_test

import (
	"grovehub/libs/core/gateway/internal/src/contracts"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tier", func() {
	It("preserves the three route values and the authorization error", func() {
		Expect(contracts.Public).To(Equal(contracts.Tier("public")))
		Expect(contracts.Protected).To(Equal(contracts.Tier("protected")))
		Expect(contracts.Restricted).To(Equal(contracts.Tier("restricted")))
		Expect(contracts.ErrWorkloadDenied).To(MatchError("workload authorization denied"))
	})
})

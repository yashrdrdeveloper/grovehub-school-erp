package contracts_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/libs/core/search/internal/src/contracts"
)

var _ = Describe("Search errors",func(){
	It("preserves domain-neutral error messages",func(){
		Expect(contracts.ErrInvalidQuery).To(MatchError("invalid search query"))
		Expect(contracts.ErrInvalidField).To(MatchError("invalid search field"))
		Expect(contracts.ErrUnsupportedOperator).To(MatchError("unsupported search operator"))
	})
})

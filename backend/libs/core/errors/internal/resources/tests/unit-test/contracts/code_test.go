package contracts_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/errors/internal/src/contracts"
)

var _ = Describe("String", func() {
	It("preserves the code exactly, including its zero value", func() {
		Expect(contracts.Code("PERSON_NOT_FOUND").String()).To(Equal("PERSON_NOT_FOUND"))
		Expect(contracts.Code("").String()).To(BeEmpty())
	})
})

package contracts_test

import (
	"grovehub/libs/core/headercontext/internal/src/contracts"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ErrInvalidHeader", func() {
	It("retains a generic error without exposing inbound values", func() {
		Expect(contracts.ErrInvalidHeader).To(MatchError("invalid request metadata header"))
	})
})

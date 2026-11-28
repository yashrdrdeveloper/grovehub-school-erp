package contracts_test

import (
	"net/http"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/errors/internal/src/contracts"
)

var _ = Describe("NewDefinition", func() {
	It("retains the validated status and declared values", func() {
		value := contracts.NewDefinition("TEST", "Public message", http.StatusConflict)

		Expect(value.Code()).To(Equal(contracts.Code("TEST")))
		Expect(value.Message()).To(Equal("Public message"))
		Expect(value.HTTPStatus()).To(Equal(http.StatusConflict))
	})
})

var _ = Describe("HTTPStatus", func() {
	It("defaults a zero-valued definition to internal server error", func() {
		Expect(contracts.Definition{}.HTTPStatus()).To(Equal(http.StatusInternalServerError))
	})
})

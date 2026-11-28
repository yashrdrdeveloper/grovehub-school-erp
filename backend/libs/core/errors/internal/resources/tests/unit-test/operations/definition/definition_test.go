package definition_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"

	"grovehub/libs/core/errors/internal/src/operations/definition"
)

var _ = Describe("NewDefinition", func() {
	It("uses the default status for legacy declarations", func() {
		value := definition.NewDefinition("LEGACY", "Legacy message")

		Expect(value.HTTPStatus()).To(Equal(http.StatusInternalServerError))
	})

	It("passes the declared HTTP error status through", func() {
		Expect(definition.NewDefinition("TEST", "Message", http.StatusUnprocessableEntity).HTTPStatus()).
			To(Equal(http.StatusUnprocessableEntity))
	})

	It("rejects an invalid declaration before constructing a definition", func() {
		Expect(func() { definition.NewDefinition("TEST", "Message", 399) }).
			To(PanicWith("error definition requires an HTTP 4xx or 5xx status"))
	})
})

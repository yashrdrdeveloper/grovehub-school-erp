package helpers_test

import (
	"grovehub/libs/core/headercontext/internal/src/contracts"
	"grovehub/libs/core/headercontext/internal/src/helpers"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReadOne", func() {
	It("accepts missing and trims a single header", func() {
		header := make(http.Header)
		value, err := helpers.ReadOne(header, "X-Request-ID")
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(BeEmpty())
		header.Set("X-Request-ID", " request-1 ")
		value, err = helpers.ReadOne(header, "X-Request-ID")
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal("request-1"))
	})

	It("rejects blank and duplicate headers even if identical", func() {
		for _, values := range [][]string{{""}, {" \t "}, {"same", "same"}, {"first", "second"}} {
			header := http.Header{"X-Request-Id": values}
			value, err := helpers.ReadOne(header, "X-Request-ID")
			Expect(err).To(BeIdenticalTo(contracts.ErrInvalidHeader))
			Expect(value).To(BeEmpty())
		}
	})
})

package errorcode_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/launcher/constants/errorcode"
)

var _ = Describe("Route error definitions", func() {
	It("keeps the route-not-found contract stable", func() {
		Expect(errorcode.RouteNotFound.Code()).To(Equal(expectedRouteNotFoundCode))
		Expect(errorcode.RouteNotFound.Message()).To(Equal(expectedRouteNotFoundMessage))
	})
})

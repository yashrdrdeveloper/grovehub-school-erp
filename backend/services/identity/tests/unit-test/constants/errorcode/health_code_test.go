package errorcode_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/services/identity/constants/errorcode"
)

var _ = Describe("Identity error definitions", func() {
	It("keeps the service-down contract stable", func() {
		Expect(errorcode.ServiceDown.Code()).To(Equal(expectedServiceDownCode))
		Expect(errorcode.ServiceDown.Message()).To(Equal(expectedServiceDownMessage))
	})
})

package constants_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/gateway/internal/src/constants"
)

var _ = Describe("Credential constants", func() {
	It("preserves the existing wire version, limits and timeouts", func() {
		Expect(constants.CredentialVersion).To(Equal("ghw1"))
		Expect(constants.MaxCredentialSize).To(Equal(8192))
		Expect(constants.MaxKeyIDLength).To(Equal(128))
		Expect(constants.MaxCredentialLifetime).To(Equal(2 * time.Minute))
		Expect(constants.OutboundCredentialLifetime).To(Equal(time.Minute))
		Expect(constants.HTTPClientTimeout).To(Equal(15 * time.Second))
	})
})

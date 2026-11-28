package constants_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/idempotency/internal/src/constants"
)

var _ = Describe("MaxKeyLength", func() {
	It("preserves the existing 128-byte limit", func() {
		Expect(constants.MaxKeyLength).To(Equal(128))
	})
})

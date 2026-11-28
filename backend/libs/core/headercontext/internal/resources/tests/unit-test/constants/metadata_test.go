package constants_test

import (
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/headercontext/internal/src/constants"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MaxMetadataIDLength", func() {
	It("preserves the existing identifier limit", func() {
		Expect(constants.MaxMetadataIDLength).To(Equal(128))
	})
})

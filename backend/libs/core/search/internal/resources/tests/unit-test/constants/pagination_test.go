package constants_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search/internal/src/constants"
)

var _ = Describe("Pagination limits",func(){
	It("preserves the shared defaults",func(){
		Expect(constants.DefaultRows).To(Equal(25))
		Expect(constants.MaxRows).To(Equal(100))
	})
})

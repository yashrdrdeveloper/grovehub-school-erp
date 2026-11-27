package logging_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	httplogging "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http/logging"
)

var _ = Describe("HTTP log level mapping", func() {
	for _, testCase := range levelCases {
		testCase := testCase
		It("maps "+testCase.name+" responses to the expected level", func() {
			Expect(httplogging.LevelForStatus(testCase.statusCode)).To(Equal(testCase.expected))
		})
	}
})

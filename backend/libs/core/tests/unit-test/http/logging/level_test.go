package logging_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corehttp "grovehub/libs/core/http"
)

var _ = Describe("HTTP log level mapping", func() {
	for _, testCase := range levelCases {
		testCase := testCase
		It("maps "+testCase.name+" responses to the expected level", func() {
			Expect(corehttp.LevelForStatus(testCase.statusCode)).To(Equal(testCase.expected))
		})
	}
})

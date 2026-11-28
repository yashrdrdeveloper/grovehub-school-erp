package constants_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/logging/internal/src/constants"
)

var _ = Describe("Standard level constants", func() {
	It("preserves standard names and count", func() {
		Expect([]string{constants.Debug, constants.Info, constants.Warn, constants.Error}).To(Equal([]string{"DEBUG", "INFO", "WARN", "ERROR"}))
		Expect(constants.StandardLevelCount).To(Equal(4))
	})
})

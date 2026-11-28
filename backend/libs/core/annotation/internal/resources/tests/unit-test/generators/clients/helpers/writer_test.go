package helpers_test

import (
	"path/filepath"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/clients/helpers"
)

var _ = Describe("WriteProxies", func() {
	It("rejects an empty discovery result without writing an output", func() {
		output := filepath.Join(GinkgoT().TempDir(),"missing.go")
		Expect(helpers.WriteProxies(nil,output)).To(MatchError(ContainSubstring("one @Client interface")))
		Expect(output).NotTo(BeAnExistingFile())
	})
})

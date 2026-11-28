package cli_test

import (
	"bytes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/libs/core/annotation/internal/src/cli"
)

var _ = Describe("Run", func() {
	It("rejects unsupported inventory options", func() {
		var out, errOut bytes.Buffer
		Expect(cli.Run([]string{"inventory", GinkgoT().TempDir(), "--unknown"}, &out, &errOut)).To(Equal(1))
		Expect(errOut.String()).To(ContainSubstring("unknown inventory option"))
		Expect(out.String()).To(BeEmpty())
	})
})

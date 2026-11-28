package cli_test

import (
	"bytes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/libs/core/annotation/internal/src/cli"
)

var _ = Describe("Run", func() {
	It("dispatches the expansion and GORM table processors rather than treating them as unknown", func() {
		var out, errOut bytes.Buffer
		root := GinkgoT().TempDir()
		Expect(cli.Run([]string{"expand", root}, &out, &errOut)).To(Equal(0))
		Expect(cli.Run([]string{"gormtable", root, root + "/out.go"}, &out, &errOut)).To(Equal(1))
		Expect(errOut.String()).To(ContainSubstring("no Go source files"))
	})
})

package cli_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"bytes"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/cli"
)

var _ = Describe("Run", func() {
	It("checks maximum positional arguments before executing a command", func() {
		var out, errOut bytes.Buffer
		Expect(cli.Run([]string{"check","first","second"},&out,&errOut)).To(Equal(2))
		Expect(errOut.String()).To(ContainSubstring("Usage: annotation check"))
		Expect(cli.Run([]string{"overlay","a","b","c"},&out,&errOut)).To(Equal(2))
	})
})

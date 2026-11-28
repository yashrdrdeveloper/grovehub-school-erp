package bootstrap_test

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/cmd/internal/src/bootstrap/impl"
)

// @Testing suite=CommandBootstrap
func CommandBootstrapSpecs() {
	Describe("CMD bootstrap", func() {
		It("dispatches an unknown command without requiring a database", func() {
			GinkgoT().Setenv("DATABASE_URL", "")
			GinkgoT().Setenv("AUDIT_DATABASE_URL", "")
			output, diagnostics := &bytes.Buffer{}, &bytes.Buffer{}
			Expect(impl.New().Build().Run([]string{"__unknown_command__"}, output, diagnostics)).To(Equal(2))
			Expect(diagnostics.String()).NotTo(BeEmpty())
		})
	})
}

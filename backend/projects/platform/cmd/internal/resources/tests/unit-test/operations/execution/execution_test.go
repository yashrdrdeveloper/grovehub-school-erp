package execution_test

import (
	"bytes"
	"io"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/cmd/internal/src/operations/execution/impl"
)

// @Testing suite=CommandExecution
func CommandExecutionSpecs() {
	Describe("CMD executor", func() {
		It("forwards arguments and streams and preserves the exit code", func() {
			output, diagnostics := &bytes.Buffer{}, &bytes.Buffer{}
			args := []string{"migrate", "all"}
			runner := impl.New(func(received []string, out, errOut io.Writer) int {
				Expect(received).To(Equal(args))
				Expect(out).To(BeIdenticalTo(output))
				Expect(errOut).To(BeIdenticalTo(diagnostics))
				_, _ = io.WriteString(out, "output")
				_, _ = io.WriteString(errOut, "diagnostic")
				return 17
			})
			Expect(runner.Run(args, output, diagnostics)).To(Equal(17))
			Expect(output.String()).To(Equal("output"))
			Expect(diagnostics.String()).To(Equal("diagnostic"))
		})
	})
}

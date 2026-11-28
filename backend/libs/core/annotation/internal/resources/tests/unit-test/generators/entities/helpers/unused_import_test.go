package helpers_test

import (
	"os"
	"path/filepath"
	"grovehub/libs/core/annotation/internal/src/generators/entities/helpers"
	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Generate", func() {
	It("omits imports not used by generated builders", func() {
		root, out := GinkgoT().TempDir(), GinkgoT().TempDir()
		fixture.Write(root, "record.go", "package entity\nimport \"time\"\nfunc now() time.Time { return time.Now() }\n// @Builder\ntype Record struct { Name string }")
		Expect(helpers.Generate(root, out)).To(Succeed())
		data, err := os.ReadFile(filepath.Join(out, "record_gombok.go"))
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).NotTo(ContainSubstring("\"time\""))
		Expect(string(data)).To(ContainSubstring("WithName"))
	})
})

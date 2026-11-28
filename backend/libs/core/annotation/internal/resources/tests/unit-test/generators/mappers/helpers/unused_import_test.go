package helpers_test

import (
	"os"
	"path/filepath"
	"grovehub/libs/core/annotation/internal/src/generators/mappers/helpers"
	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Write", func() {
	It("omits imports not used by generated mapping methods", func() {
		root := GinkgoT().TempDir()
		fixture.Write(root, "services/person/dto/output.go", "package dto\ntype Output struct { Name string }")
		fixture.Write(root, "mappers/mapper.go", "package mapper\nimport \"time\"\nimport \"example.org/demo/services/person/dto\"\nfunc now() time.Time { return time.Now() }\n// @Mapper mapperImpl\ntype Mapper interface {\n// @Mapping Name=input\nMap(input string) dto.Output\n}")
		group, found, err := helpers.Discover(root, "example.org/demo", filepath.Join(root, "mappers"))
		Expect(err).NotTo(HaveOccurred())
		Expect(found).To(BeTrue())
		out := filepath.Join(root, "generated.go")
		Expect(helpers.Write(group, out)).To(Succeed())
		data, err := os.ReadFile(out)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).NotTo(ContainSubstring("\"time\""))
		Expect(string(data)).To(ContainSubstring("dto.Output"))
	})
})

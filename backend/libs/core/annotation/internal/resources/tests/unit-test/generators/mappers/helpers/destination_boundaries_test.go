package helpers_test

import (
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
	"grovehub/libs/core/annotation/internal/src/generators/mappers/helpers"
)

var _ = Describe("Mapper destination ownership", func() {
	It("rejects an imported DTO path that is invalid or not a directory", func() {
		source := "// @Mapper impl\ntype Mapper interface{\n// @Mapping Name=input\nMap(input string) dto.Record }"
		root, dir := mapperEdgeFixture(source, "type Record struct{Name string}")
		path := filepath.Join(dir, "record.go")
		content := fixture.Read(path)
		Expect(os.WriteFile(path, []byte(strings.Replace(content, "/services/demo/dto", "/services/demo/dto?invalid", 1)), 0644)).To(Succeed())
		_, _, err := helpers.Discover(root, "example.com/project", dir)
		Expect(err).To(HaveOccurred())
		Expect(os.WriteFile(path, []byte(content), 0644)).To(Succeed())
		dtoDir := filepath.Join(root, "services", "demo", "dto")
		Expect(os.RemoveAll(dtoDir)).To(Succeed())
		Expect(os.WriteFile(dtoDir, []byte("file"), 0644)).To(Succeed())
		_, _, err = helpers.Discover(root, "example.com/project", dir)
		Expect(err).To(HaveOccurred())
	})
	It("ignores unrelated destination declarations and test source", func() {
		source := "// explanatory comment\n// @Mapper impl\ntype Mapper interface{\n// @Mapping Name=input\nMap(input string) dto.Record }"
		root, dir := mapperEdgeFixture(source, "const x=1\ntype Other struct{}\nfunc Helper(){}\ntype Record struct{Name string}")
		fixture.Write(root, "services/demo/dto/a_test.go", "invalid source")
		_, found, err := helpers.Discover(root, "example.com/project", dir)
		Expect(err).NotTo(HaveOccurred())
		Expect(found).To(BeTrue())
	})
})

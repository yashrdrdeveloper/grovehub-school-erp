package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"path/filepath"

	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
	"grovehub/libs/core/annotation/internal/src/generators/mappers/helpers"
)

var _ = Describe("Discover", func() {
	It("requires explicit mapping for every destination field", func() {
		root := GinkgoT().TempDir()
		fixture.Write(root, "services/person/dto/output.go", "package dto\ntype Output struct{ Name, Other string }")
		fixture.Write(root, "mappers/mapper.go", "package mapper\nimport \"example.org/demo/services/person/dto\"\n// @Mapper mapperImpl\ntype Mapper interface {\n// @Mapping Name=input\nMap(input string) dto.Output\n}\n")
		group, found, err := helpers.Discover(root, "example.org/demo", filepath.Join(root, "mappers"))
		Expect(err).To(MatchError(ContainSubstring("destination requires 2")))
		_ = group
		_ = found
	})
})

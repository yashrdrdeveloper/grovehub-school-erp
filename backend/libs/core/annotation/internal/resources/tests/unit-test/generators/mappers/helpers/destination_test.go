package helpers_test

import (
	"path/filepath"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/mappers/helpers"
	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
)

var _ = Describe("Discover",func(){
	It("rejects embedded destination fields",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"services/person/dto/output.go","package dto\ntype Embedded struct{}; type Output struct{ Embedded }")
		fixture.Write(root,"mappers/mapper.go","package mapper\nimport \"example.org/demo/services/person/dto\"\n// @Mapper mapperImpl\ntype Mapper interface {\n// @Mapping Name=input\nMap(input string) dto.Output\n}\n")
		group,found,err:=helpers.Discover(root,"example.org/demo",filepath.Join(root,"mappers"))
		Expect(err).To(MatchError(ContainSubstring("embedded fields are forbidden")))
		_ = group
		_ = found
	})
})

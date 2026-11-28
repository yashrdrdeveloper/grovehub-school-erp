package helpers_test

import (
	"path/filepath"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/mappers/helpers"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _ = Describe("Discover",func(){
	It("distinguishes mapper declarations from executable mapping methods",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"services/person/dto/output.go","package dto\ntype Output struct{ Name string }")
		fixture.Write(root,"mappers/mapper.go","package mapper\nimport \"example.org/demo/services/person/dto\"\n// @Mapper mapperImpl\ntype Mapper interface {\n// no mapping\nMap(input string) dto.Output\n}\n")
		group,found,err:=helpers.Discover(root,"example.org/demo",filepath.Join(root,"mappers"))
		Expect(err).NotTo(HaveOccurred())
		Expect(found).To(BeTrue())
		Expect(helpers.HasMethods(group)).To(BeFalse())
	})
})

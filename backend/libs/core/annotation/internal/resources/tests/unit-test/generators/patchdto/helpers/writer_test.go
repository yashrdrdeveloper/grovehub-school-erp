package helpers_test

import (
	"path/filepath"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/patchdto/helpers"
	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
)

var _ = Describe("Write",func(){
	It("emits conditional serialization for every present field",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"input.go","package dto\nimport core \"grovehub/libs/core/annotation\"\n// @PatchRequestDto\ntype Input struct {\nCount core.Patch[int] `json:\"count\"`\nName core.Patch[string] `json:\"name\"`\n}\n")
		specs,pkg,err:=helpers.Discover(root)
		Expect(err).NotTo(HaveOccurred())
		out:=filepath.Join(root,"output.txt")
		Expect(helpers.Write(out,pkg,specs)).To(Succeed())
		content:=fixture.Read(out)
		Expect(content).To(ContainSubstring("if request.Count.Set"))
		Expect(content).To(ContainSubstring("if request.Name.Set"))
		Expect(content).To(ContainSubstring("json.Marshal(fields)"))
	})
})

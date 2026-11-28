package helpers_test

import (
	"path/filepath"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/patchdto/helpers"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _ = Describe("Discover",func(){
	It("preserves JSON aliases instead of deriving names from Go fields",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"input.go","package dto\nimport core \"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation\"\n// @PatchRequestDto\ntype Input struct {\nDisplayName core.Patch[string] `json:\"display_name,omitempty\"`\n}\n")
		specs,pkg,err:=helpers.Discover(root)
		Expect(err).NotTo(HaveOccurred())
		Expect(specs).To(HaveLen(1))
		out:=filepath.Join(root,"output.txt")
		Expect(helpers.Write(out,pkg,specs)).To(Succeed())
		Expect(fixture.Read(out)).To(ContainSubstring(`fields["display_name"] = request.DisplayName.Value`))
	})
})

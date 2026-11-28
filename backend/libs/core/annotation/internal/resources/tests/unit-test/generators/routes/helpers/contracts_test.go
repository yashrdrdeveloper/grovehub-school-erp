package helpers_test

import (
	"path/filepath"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/routes/helpers"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _ = Describe("DiscoverWithTier",func(){
	It("supports the declared DELETE HTTP verb",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"protected/record.go","package protected\nimport \"net/http\"\n// @Routes /records\ntype RecordController interface {\n// @DELETE /{gid}\nDelete(http.ResponseWriter,*http.Request)\n}\n")
		groups,err:=helpers.DiscoverWithTier(root,"/school",true,false)
		Expect(err).NotTo(HaveOccurred())
		Expect(groups).To(HaveLen(1))
		out:=filepath.Join(root,"out.txt")
		Expect(helpers.WriteRoutes("example.org/demo","school",groups,out)).To(Succeed())
		Expect(fixture.Read(out)).To(ContainSubstring("DELETE /school/protected/records/{gid}"))
	})
})

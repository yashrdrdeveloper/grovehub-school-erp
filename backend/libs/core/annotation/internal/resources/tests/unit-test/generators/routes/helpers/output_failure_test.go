package helpers_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/routes/helpers"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("Generated output persistence",func(){
	It("preserves deterministic output and reports invalid output destinations",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"protected/source.go","package protected\n// @Routes /items\ntype ItemController interface {\n// @GET /\nRead(http.ResponseWriter,*http.Request)\n}")
		groups,err:=helpers.DiscoverWithTier(root,"/demo",true,false)
		Expect(err).NotTo(HaveOccurred())
		save:=func(path string)error{return helpers.WriteRoutes("example.com/project","demo",groups,path)}
		out:=GinkgoT().TempDir()
		path:=filepath.Join(out,"generated.go")
		Expect(save(path)).To(Succeed())
		first:=fixture.Read(path)
		Expect(first).NotTo(BeEmpty())
		Expect(save(path)).To(Succeed())
		Expect(fixture.Read(path)).To(Equal(first))
		blocked:=fixture.Write(out,"blocked","not a directory")
		Expect(save(filepath.Join(blocked,"generated.go"))).NotTo(Succeed())
		Expect(fixture.Read(blocked)).To(Equal("not a directory"))
		directory:=filepath.Join(out,"directory.go")
		Expect(os.Mkdir(directory,0755)).To(Succeed())
		Expect(save(directory)).NotTo(Succeed())
	})
})

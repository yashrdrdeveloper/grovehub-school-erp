package helpers_test

import (
	"os"
	"path/filepath"
	

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/routes/helpers"
	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("Route metadata and category boundaries",func(){
	It("rejects non-directory and cyclic controller categories",func(){
		root:=GinkgoT().TempDir()
		path:=fixture.Write(root,"protected","file")
		_,err:=helpers.DiscoverWithTier(root,"/demo",true,false)
		Expect(err).To(MatchError(ContainSubstring("controller category must be a directory")))
		Expect(os.Remove(path)).To(Succeed())
		Expect(os.Symlink("protected",path)).To(Succeed())
		_,err=helpers.DiscoverWithTier(root,"/demo",true,false)
		Expect(err).To(HaveOccurred())
	})
	It("ignores explanatory docs and test sources and rejects the wrong writer type",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"protected/ignored_test.go","invalid source")
		fixture.Write(root,"protected/controller.go","package protected\n// explanatory text\n// @Routes /records\ntype RecordController interface{\n// method explanation\n// @GET /\nRead(http.ResponseWriter,*http.Request)\n}")
		_,err:=helpers.DiscoverWithTier(root,"/demo",true,false)
		Expect(err).NotTo(HaveOccurred())
		fixture.Write(root,"protected/controller.go","package protected\n// @Routes /records\ntype RecordController interface{\n// @GET /\nRead(http.Other,*http.Request)\n}")
		_,err=helpers.DiscoverWithTier(root,"/demo",true,false)
		Expect(err).To(MatchError(ContainSubstring("must accept http.ResponseWriter")))
	})
	It("ignores unrelated service documentation and reports inaccessible prefix sources",func(){
		root:=GinkgoT().TempDir()
		path:=fixture.Write(root,"http.go","package routes\n// Other is a maintained function.\nfunc Other(){}\n// @Routes /demo\nfunc RegisterHTTP(){}")
		value,found,err:=helpers.ServicePrefix(path)
		Expect(err).NotTo(HaveOccurred());Expect(found).To(BeTrue());Expect(value).To(Equal("/demo"))
		Expect(os.Remove(path)).To(Succeed())
		Expect(os.Symlink("http.go",path)).To(Succeed())
		_,_,err=helpers.ServicePrefix(path)
		Expect(err).To(HaveOccurred())
	})
	It("propagates unreadable controller directory errors",func(){
		if os.Geteuid()==0{Skip("requires unprivileged filesystem permissions")}
		root:=GinkgoT().TempDir()
		dir:=filepath.Join(root,"protected")
		Expect(os.Mkdir(dir,0000)).To(Succeed())
		defer os.Chmod(dir,0755)
		_,err:=helpers.DiscoverWithTier(root,"/demo",true,false)
		Expect(err).To(HaveOccurred())
	})
})

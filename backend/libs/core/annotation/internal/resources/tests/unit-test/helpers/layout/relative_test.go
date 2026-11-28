package layout_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/helpers/layout"
)

var _=Describe("ResolveRelative",func(){
	It("rejects paths escaping the owning module",func(){
		root:=GinkgoT().TempDir()
		for _,path:=range []string{"","../escape","/absolute","service/../../escape"}{
			_,err:=layout.ResolveRelative(root,path)
			Expect(err).To(HaveOccurred(),path)
		}
		value,err:=layout.ResolveRelative(root,"services/demo")
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal(filepath.Join(root,"services","demo")))
	})
	It("rejects an invalid module import path in a syntactically valid directive",func(){
		root:=GinkgoT().TempDir()
		Expect(os.WriteFile(filepath.Join(root,"go.mod"),[]byte("module example.org/invalid?path\n"),0644)).To(Succeed())
		_,err:=layout.ReadModulePath(root)
		Expect(err).To(MatchError(ContainSubstring("invalid module import path")))
	})
})

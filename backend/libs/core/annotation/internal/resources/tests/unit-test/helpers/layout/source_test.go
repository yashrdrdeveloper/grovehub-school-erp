package layout_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/helpers/layout"
)

var _=Describe("GoFiles",func(){
	It("reads literal directories in order and does not interpret glob syntax",func(){
		root:=GinkgoT().TempDir()
		dir:=filepath.Join(root,"[literal]")
		Expect(os.Mkdir(dir,0755)).To(Succeed())
		for _,name:=range []string{"z.go","a.go","README.md"}{
			Expect(os.WriteFile(filepath.Join(dir,name),[]byte("package demo"),0644)).To(Succeed())
		}
		files,err:=layout.GoFiles(dir)
		Expect(err).NotTo(HaveOccurred())
		Expect(files).To(Equal([]string{filepath.Join(dir,"a.go"),filepath.Join(dir,"z.go")}))
		files,err=layout.GoFiles(filepath.Join(root,"missing"))
		Expect(err).NotTo(HaveOccurred());Expect(files).To(BeEmpty())
		_,err=layout.GoFiles(filepath.Join(dir,"a.go"))
		Expect(err).To(HaveOccurred())
	})
})

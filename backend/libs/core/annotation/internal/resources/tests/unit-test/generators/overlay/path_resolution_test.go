package overlay_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/overlay"
)

var _=Describe("Overlay path resolution",func(){
	It("rejects relative roots when the current directory was removed",func(){
		original,err:=os.Getwd();Expect(err).NotTo(HaveOccurred())
		root:=GinkgoT().TempDir()
		deleted:=filepath.Join(root,"deleted")
		Expect(os.Mkdir(deleted,0755)).To(Succeed())
		Expect(os.Chdir(deleted)).To(Succeed())
		defer func(){Expect(os.Chdir(original)).To(Succeed())}()
		Expect(os.Remove(deleted)).To(Succeed())
		Expect(overlay.WriteOverlay(".",root)).NotTo(Succeed())
		Expect(overlay.WriteOverlay(root,".")).NotTo(Succeed())
	})
	It("preserves a maintained path with a cyclic symbolic link",func(){
		root:=GinkgoT().TempDir()
		Expect(os.MkdirAll(filepath.Join(root,"pkg"),0755)).To(Succeed())
		Expect(os.MkdirAll(filepath.Join(root,".generated","pkg"),0755)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(root,".generated","pkg","zz_source.go"),[]byte("package pkg"),0644)).To(Succeed())
		path:=filepath.Join(root,"pkg","zz_source.go")
		Expect(os.Symlink("zz_source.go",path)).To(Succeed())
		Expect(overlay.WriteOverlay(root,root)).NotTo(Succeed())
		_,err:=os.Lstat(path);Expect(err).NotTo(HaveOccurred())
	})
})

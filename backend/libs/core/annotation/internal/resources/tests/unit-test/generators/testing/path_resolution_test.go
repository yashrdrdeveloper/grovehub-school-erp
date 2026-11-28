package testgen_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	testgen "grovehub/libs/core/annotation/internal/src/generators/testing"
)

var _=Describe("Test generator path resolution",func(){
	It("propagates errors from inaccessible owning module metadata",func(){
		root:=GinkgoT().TempDir()
		Expect(os.Symlink("go.mod",filepath.Join(root,"go.mod"))).To(Succeed())
		_,err:=testgen.FindModuleDir(root)
		Expect(err).To(HaveOccurred())
	})
	It("rejects relative roots when the current directory was removed",func(){
		original,err:=os.Getwd();Expect(err).NotTo(HaveOccurred())
		root:=GinkgoT().TempDir()
		deleted:=filepath.Join(root,"deleted")
		Expect(os.Mkdir(deleted,0755)).To(Succeed())
		Expect(os.Chdir(deleted)).To(Succeed())
		defer func(){Expect(os.Chdir(original)).To(Succeed())}()
		Expect(os.Remove(deleted)).To(Succeed())
		Expect(testgen.GenerateAll(".",root)).NotTo(Succeed())
		Expect(testgen.GenerateAll(root,".")).NotTo(Succeed())
	})
})

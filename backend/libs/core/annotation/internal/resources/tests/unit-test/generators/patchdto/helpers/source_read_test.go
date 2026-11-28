package helpers_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/patchdto/helpers"
)

var _=Describe("Source directory read failures",func(){
	It("propagates non-directory source paths instead of treating them as empty",func(){
		root:=GinkgoT().TempDir()
		blocked:=filepath.Join(root,"blocked")
		Expect(os.WriteFile(blocked,[]byte("not a directory"),0644)).To(Succeed())
		_,_,err:=helpers.Discover(blocked)
		Expect(err).To(HaveOccurred())
	})
})

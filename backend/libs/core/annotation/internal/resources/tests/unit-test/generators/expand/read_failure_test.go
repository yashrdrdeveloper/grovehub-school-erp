package expand_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/expand"
)

var _=Describe("Expansion source read failures",func(){
	It("rejects non-directory roots",func(){
		root:=GinkgoT().TempDir()
		path:=filepath.Join(root,"source")
		Expect(os.WriteFile(path,[]byte("file"),0644)).To(Succeed())
		Expect(expand.Expand(path)).NotTo(Succeed())
	})
})

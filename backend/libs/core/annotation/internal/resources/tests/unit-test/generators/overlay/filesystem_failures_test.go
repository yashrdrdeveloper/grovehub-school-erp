package overlay_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/overlay"
	sourceFixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("Overlay filesystem boundaries",func(){
	It("rejects generated sources from a service with no maintained package",func(){
		root,core:=GinkgoT().TempDir(),GinkgoT().TempDir()
		sourceFixture.Write(root,"services/demo/.generated/pkg/zz_source.go","package pkg")
		Expect(overlay.WriteOverlay(root,core)).To(MatchError(ContainSubstring("no original package directory")))
	})
	It("rejects looping generated roots and duplicate owner scans",func(){
		root:=GinkgoT().TempDir()
		Expect(os.Symlink(".generated",filepath.Join(root,".generated"))).To(Succeed())
		Expect(overlay.WriteOverlay(root,root)).NotTo(Succeed())
		backend:=GinkgoT().TempDir()
		owner:=filepath.Join(backend,"services","demo")
		Expect(os.MkdirAll(filepath.Join(owner,"pkg"),0755)).To(Succeed())
		sourceFixture.Write(owner,".generated/pkg/zz_source.go","package pkg")
		Expect(overlay.WriteOverlay(backend,owner)).To(MatchError(ContainSubstring("duplicate virtual file")))
	})
	It("rejects a non-directory output root without discarding source",func(){
		backend,core:=GinkgoT().TempDir(),GinkgoT().TempDir()
		Expect(os.MkdirAll(filepath.Join(backend,"services"),0755)).To(Succeed())
		Expect(os.MkdirAll(filepath.Join(core,"pkg"),0755)).To(Succeed())
		path:=sourceFixture.Write(core,".generated/pkg/zz_source.go","package pkg")
		sourceFixture.Write(backend,".generated","blocked")
		Expect(overlay.WriteOverlay(backend,core)).NotTo(Succeed())
		Expect(path).To(BeAnExistingFile())
	})
})

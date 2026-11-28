package overlay_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
	"grovehub/libs/core/annotation/internal/src/generators/overlay"
)

var _ = Describe("Overlay safety", func() {
	It("rejects missing packages, maintained source conflicts and symlinks", func() {
		root := GinkgoT().TempDir()
		generated := fixture.Write(root, ".generated/pkg/zz_source.gen.go", "package pkg")
		Expect(overlay.WriteOverlay(root, root)).To(MatchError(ContainSubstring("no original package directory")))
		logical := fixture.Write(root, "pkg/zz_source.gen.go", "package pkg")
		Expect(overlay.WriteOverlay(root, root)).To(MatchError(ContainSubstring("conflicts with maintained source")))
		Expect(os.Remove(logical)).To(Succeed())
		link := filepath.Join(root, ".generated", "pkg", "zz_link.gen.go")
		Expect(os.Symlink(generated, link)).To(Succeed())
		Expect(overlay.WriteOverlay(root, root)).To(MatchError(ContainSubstring("must be regular")))
		Expect(os.Remove(link)).To(Succeed())
		fixture.Write(root, ".generated/pkg/README.md", "ignored")
		Expect(overlay.WriteOverlay(root, root)).To(Succeed())
		out := filepath.Join(root, ".generated", "overlay.json")
		content := fixture.Read(out)
		Expect(overlay.WriteOverlay(root, root)).To(Succeed())
		Expect(fixture.Read(out)).To(Equal(content))
		Expect(os.Remove(out)).To(Succeed())
		Expect(os.Mkdir(out, 0755)).To(Succeed())
		Expect(overlay.WriteOverlay(root, root)).NotTo(Succeed())
	})
	It("rejects missing service roots and invalid generated roots", func() {
		root, core := GinkgoT().TempDir(), GinkgoT().TempDir()
		Expect(overlay.WriteOverlay(root, core)).NotTo(Succeed())
		fixture.Write(root, "services/README.md", "not a service")
		fixture.Write(core, ".generated", "not a directory")
		Expect(overlay.WriteOverlay(root, core)).NotTo(Succeed())
	})
})

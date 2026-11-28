package layout_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"

	"grovehub/libs/core/annotation/internal/src/helpers/layout"
)

var _ = Describe("ReadModulePath", func() {
	It("Read Module Path", func() {
		t := GinkgoT()
		root := t.TempDir()
		{
			_, err := layout.ReadModulePath(root)
			Expect(err).To(HaveOccurred())
		}
		filename := filepath.Join(root, "go.mod")
		for _, entry := range []struct{ source, want string }{
			{"module grovehub\n", "grovehub"},
			{"// comment\nmodule example.org/internal-platform\n", "example.org/internal-platform"},
			{"go 1.27.0\n", ""},
		} {
			Expect(os.WriteFile(filename, []byte(entry.source), 0600)).To(Succeed())
			got, err := layout.ReadModulePath(root)
			if entry.want == "" {
				Expect(err).To(HaveOccurred())
				continue
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(entry.want))
		}

	})
})

var _ = Describe("ResolveImport", func() {
	It("Resolve Import", func() {
		t := GinkgoT()
		root := t.TempDir()
		modulePath := "example.org/internal-platform"
		resolved, err := layout.ResolveImport(root, modulePath, modulePath+"/services/demo/dto")
		Expect(err).NotTo(HaveOccurred())
		Expect(resolved).To(Equal(filepath.Join(root, "services", "demo", "dto")))
		for _, invalid := range []string{"grovehub/services/demo/dto", modulePath + "/../escape", modulePath + "-other/services/demo"} {
			{
				_, err := layout.ResolveImport(root, modulePath, invalid)
				Expect(err).To(HaveOccurred())
			}
		}

	})
})

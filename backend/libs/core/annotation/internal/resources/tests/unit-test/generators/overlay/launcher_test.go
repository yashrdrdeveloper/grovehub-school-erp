package overlay_test

import (
	"encoding/json"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/libs/core/annotation/internal/src/generators/overlay"
	"os"
	"path/filepath"
)

var _ = Describe("Launcher overlay", func() {
	It("maps launcher-owned generation and rejects conflicts", func() {
		root := GinkgoT().TempDir()
		app, core := filepath.Join(root, "app"), filepath.Join(root, "core")
		owner := filepath.Join(app, "launcher")
		rel := filepath.Join("internal", "src", "routes", "zz_routes.gen.go")
		logical, generated := filepath.Join(owner, rel), filepath.Join(owner, ".generated", rel)
		for _, dir := range []string{filepath.Join(app, "services"), core, filepath.Dir(logical), filepath.Dir(generated)} {
			Expect(os.MkdirAll(dir, 0755)).To(Succeed())
		}
		Expect(os.WriteFile(generated, []byte("package routes\n"), 0600)).To(Succeed())
		Expect(overlay.WriteOverlay(app, core)).To(Succeed())
		data, err := os.ReadFile(filepath.Join(app, ".generated", "overlay.json"))
		Expect(err).NotTo(HaveOccurred())
		var config overlay.Configuration
		Expect(json.Unmarshal(data, &config)).To(Succeed())
		Expect(config.Replace).To(HaveKeyWithValue(logical, generated))
		Expect(os.WriteFile(logical, []byte("package routes\n"), 0600)).To(Succeed())
		Expect(overlay.WriteOverlay(app, core)).NotTo(Succeed())
	})
})

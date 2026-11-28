package overlay_test

import (
	"fmt"
	"encoding/json"
	"os"
	"path/filepath"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/overlay"
)

var _ = Describe("WriteOverlay", func() {
	It("Explicit Overlay Roots", func() {
		t := GinkgoT()
	root := t.TempDir()
	app := filepath.Join(root, "app")
	core := filepath.Join(root, "shared")
	owner := filepath.Join(app, "services", "demo")
	for _, dir := range []string{owner, core} {
		Expect(os.MkdirAll(dir, 0755)).To(Succeed())
	}
	Expect(overlay.WriteOverlay(core, core)).To(Succeed())
	Expect(overlay.WriteOverlay(app, core)).To(HaveOccurred())
	want := map[string]string{}
	for _, dir := range []string{owner, core} {
		logical := filepath.Join(dir, "test", "zz_example.gen.go")
		generated := filepath.Join(dir, ".generated", "test", "zz_example.gen.go")
		for _, parent := range []string{filepath.Dir(logical), filepath.Dir(generated)} {
			Expect(os.MkdirAll(parent, 0755)).To(Succeed())
		}
		Expect(os.WriteFile(generated, []byte("package test\n"), 0600)).To(Succeed())
		want[logical] = generated
	}
	Expect(overlay.WriteOverlay(app, core)).To(Succeed())
	data, err := os.ReadFile(filepath.Join(app, ".generated", "overlay.json"))
	Expect(err).NotTo(HaveOccurred())
	var config overlay.Configuration
	Expect(json.Unmarshal(data, &config)).To(Succeed())
	Expect(len(config.Replace) != len(want)).To(BeFalse(), fmt.Sprintf("got %v, want %v", config.Replace, want))
	for logical, generated := range want {
		Expect(config.Replace[logical] != generated).To(BeFalse(), fmt.Sprintf("missing mapping for %s", logical))
	}

	})
})

package overlay_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	. "github.com/onsi/ginkgo/v2"

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
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
	}
	if err := overlay.WriteOverlay(core, core); err != nil {
		t.Fatal(err)
	}
	if err := overlay.WriteOverlay(app, core); err == nil {
		t.Fatal("empty application overlay accepted")
	}
	want := map[string]string{}
	for _, dir := range []string{owner, core} {
		logical := filepath.Join(dir, "test", "zz_example.gen.go")
		generated := filepath.Join(dir, ".generated", "test", "zz_example.gen.go")
		for _, parent := range []string{filepath.Dir(logical), filepath.Dir(generated)} {
			if err := os.MkdirAll(parent, 0755); err != nil {
				t.Fatal(err)
			}
		}
		if err := os.WriteFile(generated, []byte("package test\n"), 0600); err != nil {
			t.Fatal(err)
		}
		want[logical] = generated
	}
	if err := overlay.WriteOverlay(app, core); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(app, ".generated", "overlay.json"))
	if err != nil {
		t.Fatal(err)
	}
	var config overlay.Configuration
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatal(err)
	}
	if len(config.Replace) != len(want) {
		t.Fatalf("got %v, want %v", config.Replace, want)
	}
	for logical, generated := range want {
		if config.Replace[logical] != generated {
			t.Fatalf("missing mapping for %s", logical)
		}
	}

	})
})

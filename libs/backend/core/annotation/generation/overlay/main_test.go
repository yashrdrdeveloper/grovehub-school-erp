package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestExplicitOverlayRoots(t *testing.T) {
	root := t.TempDir()
	app := filepath.Join(root, "app")
	core := filepath.Join(root, "shared")
	owner := filepath.Join(app, "services", "demo")
	for _, dir := range []string{owner, core} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
	}
	if err := writeOverlay(core, core); err != nil {
		t.Fatal(err)
	}
	if err := writeOverlay(app, core); err == nil {
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
	if err := writeOverlay(app, core); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(app, ".generated", "overlay.json"))
	if err != nil {
		t.Fatal(err)
	}
	var config configuration
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
}

package coverage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAppendManifestReportsMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing")
	if err := appendManifest(path, "generated.go"); !os.IsNotExist(err) {
		t.Fatalf("expected missing manifest error, got %v", err)
	}
}

func TestTargetRejectsInconsistentPaths(t *testing.T) {
	item := target{name: "service", owner: t.TempDir(), pattern: "."}
	if err := testTarget(".", "", "", item, "--strict"); err == nil {
		t.Fatal("expected relative path error")
	}
}

func TestTargetReportsMissingCoreDirectory(t *testing.T) {
	item := target{name: "core", owner: filepath.Join(t.TempDir(), "missing"), pattern: "./..."}
	if err := testTarget("", "", "", item, "--strict"); !os.IsNotExist(err) {
		t.Fatalf("expected missing core directory error, got %v", err)
	}
}

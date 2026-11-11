package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const patchSource = `package widget
import annotated "grovehub/libs/core/annotation/dto"

// @PatchRequestDto
type UpdateWidget struct {
	Title annotated.Patch[string] ` + "`json:\"title\"`" + `
	Count annotated.Patch[int] ` + "`json:\"count\"`" + `
}
`

func fixture(t *testing.T, source string) (string, string) {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, "services", "demo", "dto", "request", "widget")
	if err := os.MkdirAll(dir, 0755); err != nil { t.Fatal(err) }
	file := filepath.Join(dir, "widget.go")
	if err := os.WriteFile(file, []byte(source), 0644); err != nil { t.Fatal(err) }
	return root, filepath.Join(root, "services", "demo", ".generated", "dto", "request", "widget", "zz_patchdto.gen.go")
}

func TestGeneratePatchRequestDTO(t *testing.T) {
	root, output := fixture(t, patchSource)
	if err := generateAll(root); err != nil { t.Fatal(err) }
	generated, err := os.ReadFile(output)
	if err != nil { t.Fatal(err) }
	for _, expected := range []string{
		"package widget",
		"func (request UpdateWidget) MarshalJSON() ([]byte, error)",
		"if request.Title.Set {",
		"fields[\"title\"] = request.Title.Value",
		"fields[\"count\"] = request.Count.Value",
	} {
		if !strings.Contains(string(generated), expected) { t.Fatalf("missing %q: %s", expected, generated) }
	}
	if err := generateAll(root); err != nil { t.Fatal(err) }
	again, err := os.ReadFile(output)
	if err != nil || string(again) != string(generated) { t.Fatalf("generator not idempotent: %v", err) }
}

func TestRejectInvalidPatchAnnotations(t *testing.T) {
	cases := map[string]string{
		"wrong type": `package widget
// @PatchRequestDto
type UpdateWidget struct { Name string ` + "`json:\"name\"`" + ` }`,
		"missing json tag": `package widget
import dto "grovehub/libs/core/annotation/dto"
// @PatchRequestDto
type UpdateWidget struct { Name dto.Patch[string] }`,
		"manual MarshalJSON": patchSource + `
func (UpdateWidget) MarshalJSON() ([]byte, error) { return nil, nil }`,
		"bad json tag": strings.Replace(patchSource, "` + "`json:\"count\"`" + `", "` + "`json:\"title\"`" + `", 1),
	}
	for name, src := range cases {
		t.Run(name, func(t *testing.T) {
			root, _ := fixture(t, src)
			if err := generateAll(root); err == nil { t.Fatal("expected generator to reject invalid PATCH DTO") }
		})
	}
}

func TestRemoveStaleGeneratedPatchMethod(t *testing.T) {
	root, output := fixture(t, patchSource)
	if err := generateAll(root); err != nil { t.Fatal(err) }
	src := strings.Replace(patchSource, "// @PatchRequestDto", "// no annotation", 1)
	path := filepath.Join(root, "services", "demo", "dto", "request", "widget", "widget.go")
	if err := os.WriteFile(path, []byte(src), 0644); err != nil { t.Fatal(err) }
	if err := generateAll(root); err != nil { t.Fatal(err) }
	if _, err := os.Stat(output); !os.IsNotExist(err) { t.Fatalf("stale generated output retained: %v", err) }
}

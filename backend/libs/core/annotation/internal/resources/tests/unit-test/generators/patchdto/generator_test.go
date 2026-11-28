package patchdto_test

import (
	"os"
	"path/filepath"
	"strings"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/patchdto"
)

const patchSource = `package widget
import annotated "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation"

// @PatchRequestDto
type UpdateWidget struct {
	Title annotated.Patch[string] ` + "`json:\"title\"`" + `
	Count annotated.Patch[int] ` + "`json:\"count\"`" + `
}
`

func fixture(t GinkgoTInterface, source string) (string, string) {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, "services", "demo", "dto", "request", "widget")
	if err := os.MkdirAll(dir, 0755); err != nil { t.Fatal(err) }
	file := filepath.Join(dir, "widget.go")
	if err := os.WriteFile(file, []byte(source), 0644); err != nil { t.Fatal(err) }
	return root, filepath.Join(root, "services", "demo", ".generated", "dto", "request", "widget", "zz_patchdto.gen.go")
}

var _ = Describe("GenerateAll", func() {
	It("Generate Patch Request DTO", func() {
		t := GinkgoT()
	root, output := fixture(t, patchSource)
	if err := patchdto.GenerateAll(root); err != nil { t.Fatal(err) }
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
	if err := patchdto.GenerateAll(root); err != nil { t.Fatal(err) }
	again, err := os.ReadFile(output)
	if err != nil || string(again) != string(generated) { t.Fatalf("generator not idempotent: %v", err) }

	})
})

var _ = Describe("GenerateAll", func() {
	It("Reject Invalid Patch Annotations", func() {
		t := GinkgoT()
	cases := map[string]string{
		"wrong type": `package widget
// @PatchRequestDto
type UpdateWidget struct { Name string ` + "`json:\"name\"`" + ` }`,
		"missing json tag": `package widget
import dto "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation"
// @PatchRequestDto
type UpdateWidget struct { Name dto.Patch[string] }`,
		"manual MarshalJSON": patchSource + `
func (UpdateWidget) MarshalJSON() ([]byte, error) { return nil, nil }`,
		"bad json tag": strings.Replace(patchSource, "`json:\"count\"`", "`json:\"title\"`", 1),
	}
	for name, src := range cases {
		By(name)
		{
			root, _ := fixture(t, src)
			if err := patchdto.GenerateAll(root); err == nil { t.Fatal("expected generator to reject invalid PATCH DTO") }
		}
	}

	})
})

var _ = Describe("GenerateAll", func() {
	It("Remove Stale Generated Patch Method", func() {
		t := GinkgoT()
	root, output := fixture(t, patchSource)
	if err := patchdto.GenerateAll(root); err != nil { t.Fatal(err) }
	src := strings.Replace(patchSource, "// @PatchRequestDto", "// no annotation", 1)
	path := filepath.Join(root, "services", "demo", "dto", "request", "widget", "widget.go")
	if err := os.WriteFile(path, []byte(src), 0644); err != nil { t.Fatal(err) }
	if err := patchdto.GenerateAll(root); err != nil { t.Fatal(err) }
	if _, err := os.Stat(output); !os.IsNotExist(err) { t.Fatalf("stale generated output retained: %v", err) }

	})
})

var _ = Describe("GenerateAll", func() {
	It("preserves legacy PATCH imports during facade migration", func() {
		root, output := fixture(GinkgoT(), strings.Replace(patchSource, "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation", "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/dto", 1))
		Expect(patchdto.GenerateAll(root)).To(Succeed())
		generated, err := os.ReadFile(output)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(generated)).To(ContainSubstring("func (request UpdateWidget) MarshalJSON()"))
	})
})

package patchdto_test

import (
	"fmt"
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
	Expect(os.MkdirAll(dir, 0755)).To(Succeed())
	file := filepath.Join(dir, "widget.go")
	Expect(os.WriteFile(file, []byte(source), 0644)).To(Succeed())
	return root, filepath.Join(root, "services", "demo", ".generated", "dto", "request", "widget", "zz_patchdto.gen.go")
}

var _ = Describe("GenerateAll", func() {
	It("Generate Patch Request DTO", func() {
		t := GinkgoT()
	root, output := fixture(t, patchSource)
	Expect(patchdto.GenerateAll(root)).To(Succeed())
	generated, err := os.ReadFile(output)
	Expect(err).NotTo(HaveOccurred())
	for _, expected := range []string{
		"package widget",
		"func (request UpdateWidget) MarshalJSON() ([]byte, error)",
		"if request.Title.Set {",
		"fields[\"title\"] = request.Title.Value",
		"fields[\"count\"] = request.Count.Value",
	} {
		Expect(string(generated)).To(ContainSubstring(expected))
	}
	Expect(patchdto.GenerateAll(root)).To(Succeed())
	again, err := os.ReadFile(output)
	Expect(err != nil || string(again) != string(generated)).To(BeFalse(), fmt.Sprintf("generator not idempotent: %v", err))

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
			Expect(patchdto.GenerateAll(root)).To(HaveOccurred())
		}
	}

	})
})

var _ = Describe("GenerateAll", func() {
	It("Remove Stale Generated Patch Method", func() {
		t := GinkgoT()
	root, output := fixture(t, patchSource)
	Expect(patchdto.GenerateAll(root)).To(Succeed())
	src := strings.Replace(patchSource, "// @PatchRequestDto", "// no annotation", 1)
	path := filepath.Join(root, "services", "demo", "dto", "request", "widget", "widget.go")
	Expect(os.WriteFile(path, []byte(src), 0644)).To(Succeed())
	Expect(patchdto.GenerateAll(root)).To(Succeed())
	{
		_, err := os.Stat(output)
		Expect(os.IsNotExist(err)).To(BeTrue())
	}

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

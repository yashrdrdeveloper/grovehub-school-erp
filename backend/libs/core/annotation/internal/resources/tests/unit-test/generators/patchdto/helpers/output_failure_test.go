package helpers_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
	"grovehub/libs/core/annotation/internal/src/generators/patchdto/helpers"
)

var _ = Describe("Generated output persistence", func() {
	It("preserves deterministic output and reports invalid output destinations", func() {
		root := GinkgoT().TempDir()
		fixture.Write(root, "source.go", "package dto\nimport annotation \"grovehub/libs/core/annotation\"\n// @PatchRequestDto\ntype Input struct{ Name annotation.Patch[string] `json:\"name\"` }")
		specs, pkg, err := helpers.Discover(root)
		Expect(err).NotTo(HaveOccurred())
		save := func(path string) error { return helpers.Write(path, pkg, specs) }
		out := GinkgoT().TempDir()
		path := filepath.Join(out, "generated.go")
		Expect(save(path)).To(Succeed())
		first := fixture.Read(path)
		Expect(first).NotTo(BeEmpty())
		Expect(save(path)).To(Succeed())
		Expect(fixture.Read(path)).To(Equal(first))
		blocked := fixture.Write(out, "blocked", "not a directory")
		Expect(save(filepath.Join(blocked, "generated.go"))).NotTo(Succeed())
		Expect(fixture.Read(blocked)).To(Equal("not a directory"))
		directory := filepath.Join(out, "directory.go")
		Expect(os.Mkdir(directory, 0755)).To(Succeed())
		Expect(save(directory)).NotTo(Succeed())
	})
})

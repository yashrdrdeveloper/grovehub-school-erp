package patchdto_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sourceFixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
	"grovehub/libs/core/annotation/internal/src/generators/patchdto"
)

var _ = Describe("Patch DTO overlay lifecycle", func() {
	It("ignores generated and vendor trees and removes stale outputs safely", func() {
		root := GinkgoT().TempDir()
		owner := filepath.Join(root, "services", "demo")
		sourceFixture.Write(owner, "internal/app/dto/update.go", "package dto\nimport annotation \"grovehub/libs/core/annotation\"\n// @PatchRequestDto\ntype Update struct{ Name annotation.Patch[string] `json:\"name\"` }")
		sourceFixture.Write(owner, "internal/app/dto/vendor/broken.go", "invalid source")
		sourceFixture.Write(owner, "internal/app/dto/.generated/broken.go", "invalid source")
		Expect(patchdto.GenerateAll(root)).To(Succeed())
		out := filepath.Join(owner, ".generated", "internal", "app", "dto", "zz_patchdto.gen.go")
		Expect(sourceFixture.Read(out)).To(ContainSubstring("MarshalJSON"))
		Expect(os.Remove(out)).To(Succeed())
		Expect(os.MkdirAll(filepath.Join(out, "child"), 0755)).To(Succeed())
		Expect(patchdto.GenerateAll(root)).NotTo(Succeed())
		sourceFixture.Write(owner, "internal/app/dto/update.go", "package dto\ntype Update struct{}")
		Expect(patchdto.GenerateAll(root)).NotTo(Succeed())
		Expect(os.RemoveAll(out)).To(Succeed())
		Expect(patchdto.GenerateAll(root)).To(Succeed())
		Expect(out).NotTo(BeAnExistingFile())
		Expect(os.Symlink("client", filepath.Join(owner, "client"))).To(Succeed())
		Expect(patchdto.GenerateAll(root)).NotTo(Succeed())
	})
})

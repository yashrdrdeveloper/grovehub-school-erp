package helpers_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
	contract "grovehub/libs/core/annotation/internal/src/contracts"
	"grovehub/libs/core/annotation/internal/src/generators/check/helpers"
)

var _ = Describe("Validation annotation checking", func() {
	It("validates nested fields and rejects invalid source or annotation syntax", func() {
		root := GinkgoT().TempDir()
		Expect(helpers.CheckValidationTags(filepath.Join(root, "missing.go"))).NotTo(Succeed())
		path := fixture.Write(root, "source.go", "package dto\ntype Input struct {\n// explanatory comment\n// @Valid\nNested struct{} `validate:\"valid\"`\n// @NotBlank\nName string `validate:\"notblank\"`\n// @NotNull\nOptional *string `validate:\"nonnull\"`\n// @Sensitive\nSecret string\n}")
		Expect(helpers.CheckValidationTags(path)).To(Succeed())
		for _, annotation := range []string{"@Size", "@Size max=bad", "@Size max=2 extra", "@UnknownDirective", "@NotBlank"} {
			fixture.Write(root, "source.go", "package dto\ntype Input struct{\n// "+annotation+"\nName string\nUnrelated string\n}")
			Expect(helpers.CheckValidationTags(path)).NotTo(Succeed(), annotation)
		}
	})
	It("ignores non-provider annotations and reports malformed provider entries", func() {
		index := make(helpers.ProviderIndex)
		Expect(index.Add([]contract.Occurrence{{Directive: contract.Directive{Name: "Entity", Arguments: "demo.records"}}})).To(Succeed())
		Expect(index).To(BeEmpty())
		Expect(index.Add([]contract.Occurrence{{Directive: contract.Directive{Name: "Provider", Arguments: "bad"}, File: "provider.go", Line: 5}})).To(MatchError(ContainSubstring("provider.go:5")))
	})
})

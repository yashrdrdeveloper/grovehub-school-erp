package registry_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/registry"
)

var _ = Describe("BuildDefinitions", func() {
	It("rejects incomplete catalog rows and mismatched syntax", func() {
		Expect(func() { registry.BuildDefinitions([]registry.UseCase{{Name: "Sample"}}, nil) }).To(Panic())
		valid := registry.UseCase{Name: "Sample", Syntax: "@Other", Category: "Test", Purpose: "Test catalog validation", Target: registry.TypeTarget}
		Expect(func() { registry.BuildDefinitions([]registry.UseCase{valid}, nil) }).To(Panic())
		valid.Syntax = "@Sample"
		defs := registry.BuildDefinitions([]registry.UseCase{valid}, nil)
		Expect(defs).To(HaveKey("Sample"))
		Expect(defs["Sample"].UseCases).To(Equal([]registry.UseCase{valid}))
	})
	It("rejects empty schema and table identifiers without accepting duplicate fields", func() {
		for _, source := range []string{"// @Entity .records", "// @Entity schema.", "// @Filterable fields=gid,gid", "// @Sortable fields=gid,"} {
			_, _, err := registry.ParseComment(source)
			Expect(err).To(HaveOccurred(), source)
		}
	})
})

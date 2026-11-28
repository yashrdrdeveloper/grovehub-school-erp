package registry_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/registry"
)

var _ = Describe("Annotation registry", func() {
	It("provides deterministic, independently owned catalog metadata", func() {
		all := registry.All()
		Expect(all).NotTo(BeEmpty())
		for i := 1; i < len(all); i++ {
			Expect(all[i].Name > all[i-1].Name).To(BeTrue())
		}
		name := all[0].Name
		original := all[0].UseCases[0].Name
		all[0].UseCases[0].Name = "mutated"
		entry, found := registry.Lookup("@" + name)
		Expect(found).To(BeTrue())
		Expect(entry.UseCases[0].Name).To(Equal(original))
		uses := registry.Uses()
		Expect(uses).NotTo(BeEmpty())
		uses[0].Name = "changed"
		Expect(registry.Uses()[0].Name).NotTo(Equal("changed"))
		patterns := registry.Patterns()
		Expect(patterns).To(HaveLen(23))
		patterns[0].Name = "changed"
		Expect(registry.Patterns()[0].Name).NotTo(Equal("changed"))
		_, found = registry.Lookup("does-not-exist")
		Expect(found).To(BeFalse())
	})
	It("accepts supported syntax and rejects incomplete or invalid options", func() {
		for _, line := range []string{"// @Entity demo.records mode=table", "// @Repository entity=Record", "// @Searchable config=cfg alias=p", "// @Filterable fields=gid,name", "// @Sortable fields=createdAt", "// @Cursor field=person.gid"} {
			d, found, err := registry.ParseComment(line)
			Expect(err).NotTo(HaveOccurred(), line)
			Expect(found).To(BeTrue())
			Expect(registry.ValidateDirective(d)).To(Succeed())
		}
		for _, line := range []string{"// @Entity", "// @Entity Demo.records", "// @Entity demo.bad-name", "// @Repository", "// @Repository wrong=Record", "// @Repository entity=", "// @Repository entity=A entity=B", "// @Repository entity=A-B", "// @Searchable config=cfg", "// @Searchable alias=p", "// @Searchable config=cfg alias=p impl=9bad", "// @Filterable fields=", "// @Filterable fields=gid,", "// @Cursor field=foo..bar"} {
			_, found, err := registry.ParseComment(line)
			Expect(found).To(BeTrue(), line)
			Expect(err).To(MatchError(ContainSubstring("invalid annotation directive")), line)
		}
		_, found, err := registry.ParseComment("ordinary text")
		Expect(err).NotTo(HaveOccurred())
		Expect(found).To(BeFalse())
		Expect(registry.ValidateDirective(registry.Directive{Name: "Missing"})).To(MatchError(ContainSubstring("unknown annotation directive")))
	})
})

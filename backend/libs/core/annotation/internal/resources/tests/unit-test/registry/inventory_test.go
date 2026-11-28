package registry_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sort"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/registry"
)

var _ = Describe("Uses", func() {
	It("preserves catalog rows while All returns unique sorted definitions", func() {
		uses:=registry.Uses()
		Expect(uses).NotTo(BeEmpty())
		for _,entry:=range uses {
			definition,found:=registry.Lookup(entry.Name)
			Expect(found).To(BeTrue())
			Expect(definition.UseCases).To(ContainElement(entry))
		}
		first:=uses[0]
		uses[0].Name="changed"
		Expect(registry.Uses()[0]).To(Equal(first))
		var names []string
		seen:=map[string]bool{}
		for _,definition:=range registry.All(){
			Expect(seen[definition.Name]).To(BeFalse())
			seen[definition.Name]=true
			names=append(names,definition.Name)
		}
		Expect(sort.StringsAreSorted(names)).To(BeTrue())
	})
})

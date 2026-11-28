package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go/parser"
	"go/token"
	"grovehub/libs/core/annotation/internal/src/parser/helpers"
	"grovehub/libs/core/annotation/internal/src/registry"
)

var _ = Describe("CollectOccurrences", func() {
	It("keeps source positions, skips ordinary comments and rejects invalid directives", func() {
		set := token.NewFileSet()
		file, err := parser.ParseFile(set, "demo.go", "package demo\n// ordinary\n// @Entity demo.records\ntype Record struct{}\n", parser.ParseComments)
		Expect(err).NotTo(HaveOccurred())
		items, err := helpers.CollectOccurrences(set, file.Comments, helpers.CollectTargets(file))
		Expect(err).NotTo(HaveOccurred())
		Expect(items).To(HaveLen(1))
		Expect(items[0].Line).To(Equal(3))
		Expect(items[0].Column).To(Equal(1))
		Expect(items[0].Directive.Arguments).To(Equal("demo.records"))
		set = token.NewFileSet()
		file, err = parser.ParseFile(set, "bad.go", "package demo\n// @UnknownDirective\ntype Record struct{}\n", parser.ParseComments)
		Expect(err).NotTo(HaveOccurred())
		items, err = helpers.CollectOccurrences(set, file.Comments, helpers.CollectTargets(file))
		Expect(err).To(MatchError(registry.ErrUnknownDirective))
		Expect(items).To(BeNil())
	})
})

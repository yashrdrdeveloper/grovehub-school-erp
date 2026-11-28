package annotation_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation"
)

var _ = Describe("Lookup", func() {
	It("returns caller-owned metadata", func() {
		definition, found := annotation.Lookup("Entity")
		Expect(found).To(BeTrue())
		Expect(definition.UseCases).NotTo(BeEmpty())
		definition.UseCases[0].Name = "modified"

		next, found := annotation.Lookup("Entity")
		Expect(found).To(BeTrue())
		Expect(next.UseCases[0].Name).To(Equal("Entity"))
	})
})

var _ = Describe("ParseComment", func() {
	It("validates recognized and unknown directives through the facade", func() {
		directive, found, err := annotation.ParseComment("// @Entity sample.records")
		Expect(err).NotTo(HaveOccurred())
		Expect(found).To(BeTrue())
		Expect(directive.Name).To(Equal("Entity"))
		Expect(directive.Arguments).To(Equal("sample.records"))

		_, _, err = annotation.ParseComment("// @UnknownAnnotation")
		Expect(err).To(MatchError(annotation.ErrUnknownDirective))
	})
})

var _ = Describe("ScanFile", func() {
	It("identifies annotated declarations from Go source", func() {
		path := filepath.Join(GinkgoT().TempDir(), "sample.go")
		source := "package sample\n// @Entity sample.records\ntype Record struct{}\n"
		Expect(os.WriteFile(path, []byte(source), 0600)).To(Succeed())

		occurrences, err := annotation.ScanFile(path)
		Expect(err).NotTo(HaveOccurred())
		Expect(occurrences).To(HaveLen(1))
		Expect(occurrences[0].Directive.Name).To(Equal("Entity"))
		Expect(occurrences[0].Target).To(Equal(annotation.TypeTarget))
	})
})

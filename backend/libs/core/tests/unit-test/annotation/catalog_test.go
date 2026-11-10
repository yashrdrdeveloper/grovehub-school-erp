package annotation_test

import (
    "errors"
    "os"
    "path/filepath"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    "grovehub/libs/core/annotation"
    sourceparser "grovehub/libs/core/annotation/generation/parser"
)

var _ = Describe("Phase 1 annotation catalog", func() {
    It("defines every documented use case and all 23 design patterns", func() {
        Expect(annotation.Uses()).To(HaveLen(226))
        Expect(annotation.Patterns()).To(HaveLen(23))
        Expect(annotation.All()).To(HaveLen(232))
        for _, name := range []string{
            "Entity", "Mapper", "Validator", "Valid", "NonNull", "Repository",
            "Service", "Client", "Routes", "Event", "Producer", "Consumer",
            "Configuration", "Log", "Trace", "Server", "Mock", "Strategy",
            "Bridge", "Schema", "UUID", "MappingIgnore",
        } {
            _, ok := annotation.Lookup(name)
            Expect(ok).To(BeTrue(), name)
        }
    })

    It("returns immutable-to-caller catalog values", func() {
        first, ok := annotation.Lookup("@Entity")
        Expect(ok).To(BeTrue())
        first.UseCases[0].Purpose = "changed"
        again, _ := annotation.Lookup("Entity")
        Expect(again.UseCases[0].Purpose).NotTo(Equal("changed"))
    })

    It("parses known directives without activating any runtime behavior", func() {
        d, found, err := annotation.ParseComment("// @Entity foundation.person_contacts")
        Expect(err).NotTo(HaveOccurred())
        Expect(found).To(BeTrue())
        Expect(d.Name).To(Equal("Entity"))
        Expect(d.Arguments).To(Equal("foundation.person_contacts"))

        d, found, err = annotation.ParseComment("// @NonNull identity.InvalidRegistration")
        Expect(err).NotTo(HaveOccurred())
        Expect(found).To(BeTrue())
        Expect(d.Arguments).To(Equal("identity.InvalidRegistration"))

        _, found, err = annotation.ParseComment("// plain comment")
        Expect(found).To(BeFalse())
        Expect(err).NotTo(HaveOccurred())
    })

    It("rejects unknown names and invalid entity tables", func() {
        _, _, err := annotation.ParseComment("// @MadeUpAction")
        Expect(errors.Is(err, annotation.ErrUnknownDirective)).To(BeTrue())
        _, _, err = annotation.ParseComment("// @Entity invalid")
        Expect(errors.Is(err, annotation.ErrInvalidDirective)).To(BeTrue())
        _, _, err = annotation.ParseComment("// @Entity foundation.person_contacts;DROP")
        Expect(errors.Is(err, annotation.ErrInvalidDirective)).To(BeTrue())
    })

    It("discovers source metadata without modifying service code", func() {
        root := GinkgoT().TempDir()
        filename := filepath.Join(root, "sample.go")
        contents := "package sample\n// @Entity foundation.people\ntype Person struct {\n // @NonNull person.InvalidName\n Name *string\n}\n"
        Expect(os.WriteFile(filename, []byte(contents), 0600)).To(Succeed())
        occurrences, err := sourceparser.ScanFile(filename)
        Expect(err).NotTo(HaveOccurred())
        Expect(occurrences).To(HaveLen(2))
        Expect(occurrences[0].Directive.Name).To(Equal("Entity"))
        Expect(occurrences[0].Target).To(Equal(annotation.TypeTarget))
        Expect(occurrences[1].Directive.Name).To(Equal("NonNull"))
        Expect(occurrences[1].Target).To(Equal(annotation.FieldTarget))
    })
})

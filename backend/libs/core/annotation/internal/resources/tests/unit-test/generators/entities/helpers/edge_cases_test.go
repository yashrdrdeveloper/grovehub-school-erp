package helpers_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
	"grovehub/libs/core/annotation/internal/src/generators/entities/helpers"
)

var _ = Describe("Entity accessor generation edge cases", func() {
	It("preserves manual accessors, ignores generated sources, removes stale output and keeps unrelated files", func() {
		root, out := GinkgoT().TempDir(), GinkgoT().TempDir()
		fixture.Write(root, "ignored_test.go", "invalid source")
		fixture.Write(root, "ignored_gombok.go", "invalid source")
		fixture.Write(root, "ignored.gen.go", "invalid source")
		fixture.Write(root, "record.go", "package entity\nimport \"time\"\n// @Builder\n// @Getter\n// @Setter\ntype Record struct{ Created time.Time; Name string }\nfunc (r *Record) GetName() string{return r.Name}\nfunc (r *Record) SetName(name string){r.Name=name}\n")
		fixture.Write(root, "plain.go", "package entity\ntype Plain struct{}\ntype Alias string\nfunc Helper(){}")
		stale := fixture.Write(out, "removed_gombok.go", "stale")
		untouched := fixture.Write(out, "zz_tables.gen.go", "unrelated")
		Expect(helpers.Generate(root, out)).To(Succeed())
		output := filepath.Join(out, "record_gombok.go")
		content := fixture.Read(output)
		Expect(content).To(ContainSubstring("NewRecordBuilder"))
		Expect(content).To(ContainSubstring("WithCreated"))
		Expect(content).To(ContainSubstring("\"time\""))
		Expect(content).NotTo(ContainSubstring("func (value *Record) Get"))
		Expect(content).NotTo(ContainSubstring("func (value *Record) Set"))
		Expect(stale).NotTo(BeAnExistingFile())
		Expect(untouched).To(BeAnExistingFile())
		Expect(helpers.Generate(root, out)).To(Succeed())
		Expect(fixture.Read(output)).To(Equal(content))
		fixture.Write(root, "record.go", "package entity\ntype Record struct{}")
		Expect(helpers.Generate(root, out)).To(Succeed())
		Expect(output).NotTo(BeAnExistingFile())
	})
	It("rejects conflicting entity directives, embedded fields, mixed packages and malformed source", func() {
		for _, test := range []struct{ source, reason string }{
			{"// @Entity demo.records\n// @Builder\ntype Record struct{}", "cannot request builders"},
			{"// @Builder\ntype Record struct{ Embedded }", "embedded fields"},
			{"func broken(", "expected"},
		} {
			root := GinkgoT().TempDir()
			fixture.Write(root, "record.go", "package entity\n"+test.source)
			Expect(helpers.Generate(root, GinkgoT().TempDir())).To(MatchError(ContainSubstring(test.reason)))
		}
		root := GinkgoT().TempDir()
		fixture.Write(root, "a.go", "package a")
		fixture.Write(root, "b.go", "package b")
		Expect(helpers.Generate(root, GinkgoT().TempDir())).To(MatchError(ContainSubstring("mixed packages")))
	})
	It("reports output path failures without overwriting maintained files", func() {
		root, out := GinkgoT().TempDir(), GinkgoT().TempDir()
		fixture.Write(root, "record.go", "package entity\n// @Setter\ntype Record struct{ Name string }")
		blocked := fixture.Write(out, "blocked", "not a directory")
		Expect(helpers.Generate(root, blocked)).NotTo(Succeed())
		Expect(os.Mkdir(filepath.Join(out, "record_gombok.go"), 0755)).To(Succeed())
		Expect(helpers.Generate(root, out)).NotTo(Succeed())
		Expect(fixture.Read(blocked)).To(Equal("not a directory"))
	})
})

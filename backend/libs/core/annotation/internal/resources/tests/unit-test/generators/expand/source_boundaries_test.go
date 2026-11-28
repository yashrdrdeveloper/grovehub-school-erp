package expand_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
	"grovehub/libs/core/annotation/internal/src/generators/expand"
)

var _ = Describe("Expansion source boundaries", func() {
	It("ignores generic receivers, ordinary declarations and unannotated structures", func() {
		root := GinkgoT().TempDir()
		path := fixture.Write(root, "record.go", "package entity\ntype Generic[T any] struct{}\nfunc (r Generic[T]) Get() T {var v T;return v}\ntype Plain struct{}\n// ordinary docs\ntype Other struct{}\n")
		before := fixture.Read(path)
		Expect(expand.Expand(root)).To(Succeed())
		Expect(fixture.Read(path)).To(Equal(before))
		Expect(os.Mkdir(filepath.Join(root, "directory.go"), 0755)).To(Succeed())
		Expect(expand.Expand(root)).NotTo(Succeed())
	})
	It("returns a write failure when an annotated staging source cannot be changed", func() {
		if os.Geteuid() == 0 {
			Skip("requires unprivileged filesystem permissions")
		}
		root := GinkgoT().TempDir()
		path := fixture.Write(root, "record.go", "package entity\n// @Entity demo.record\ntype Record struct{}")
		before := fixture.Read(path)
		Expect(os.Chmod(path, 0444)).To(Succeed())
		defer os.Chmod(path, 0644)
		Expect(expand.Expand(root)).NotTo(Succeed())
		Expect(fixture.Read(path)).To(Equal(before))
	})
})

package coverage_test

import (
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/testing/coverage"
)

var _ = Describe("Coverage runner failure handling", func() {
	It("rejects missing service roots and mixed all selections, and deduplicates names", func() {
		root := GinkgoT().TempDir()
		_, err := coverage.Resolve(root, root, "all", "")
		Expect(err).To(HaveOccurred())
		for _, dir := range []string{"launcher", "services/demo", "services/.hidden"} {
			Expect(os.MkdirAll(filepath.Join(root, dir), 0755)).To(Succeed())
		}
		writeFixture(filepath.Join(root, "services", "README"), "not a service")
		_, err = coverage.Resolve(root, root, "", "all,demo")
		Expect(err).To(MatchError("all cannot be combined with selected packages"))
		items, err := coverage.Resolve(root, root, "", "demo,demo,core")
		Expect(err).NotTo(HaveOccurred())
		Expect(coverage.Names(items)).To(Equal("core, demo"))
		items, err = coverage.Resolve(root, root, "all", "")
		Expect(err).NotTo(HaveOccurred())
		Expect(coverage.Names(items)).To(Equal("core, demo, launcher"))
	})
	It("preserves unsafe stale generated manifests and rejects blocked output directories", func() {
		root, _, _, args := coverageFixture()
		generated := filepath.Join(root, ".generated")
		writeFixture(generated, "not a directory")
		Expect(coverage.Run(args)).NotTo(Succeed())
		Expect(os.Remove(generated)).To(Succeed())
		manifest := filepath.Join(generated, "unit-coverage", ".materialized-files")
		outside := filepath.Join(GinkgoT().TempDir(), "zz_outside.go")
		writeFixture(outside, "// Code generated.\npackage example")
		writeFixture(manifest, outside+"\n")
		Expect(coverage.Run(args)).To(MatchError(ContainSubstring("unsafe stale generated path")))
		Expect(outside).To(BeAnExistingFile())
		Expect(manifest).To(BeAnExistingFile())
	})
	It("cleans partial materialization when an output is blocked or a source cannot be copied", func() {
		root := GinkgoT().TempDir()
		overlay, manifest := filepath.Join(root, "overlay.json"), filepath.Join(root, "manifest")
		source := filepath.Join(root, "source.go")
		writeFixture(source, "package example")
		blocked := filepath.Join(root, "blocked")
		writeFixture(blocked, "not a directory")
		writeOverlay(overlay, map[string]string{filepath.Join(blocked, "child.go"): source})
		_, err := coverage.Materialize(overlay, manifest)
		Expect(err).To(HaveOccurred())
		Expect(manifest).NotTo(BeAnExistingFile())
		target := filepath.Join(root, "zz_generated.go")
		// Opening a directory as a source succeeds on Unix, but reading it does not.
		writeOverlay(overlay, map[string]string{target: root})
		_, err = coverage.Materialize(overlay, manifest)
		Expect(err).To(HaveOccurred())
		Expect(target).NotTo(BeAnExistingFile())
		Expect(manifest).NotTo(BeAnExistingFile())
	})
	It("rejects oversized coverage records instead of accepting a truncated profile", func() {
		path := filepath.Join(GinkgoT().TempDir(), "coverage.out")
		writeFixture(path, "mode: atomic\n"+strings.Repeat("x", 128*1024)+" 1 0\n")
		_, _, err := coverage.Score(path)
		Expect(err).To(MatchError(ContainSubstring("token too long")))
	})
	It("rejects invalid production packages and failed unit tests", func() {
		root, _, _, args := coverageFixture()
		source := filepath.Join(root, "services", "sample", "sample.go")
		writeFixture(source, "package sample\nfunc broken(")
		Expect(coverage.Run(args)).NotTo(Succeed())
		writeFixture(source, "package sample\nfunc Answer() int { return 7 }")
		Expect(coverage.Run(args)).To(MatchError("one or more selected unit coverage runs failed"))
	})
	It("counts a core package without tests and rejects empty production coverage", func() {
		root, generator, reports, _ := coverageFixture()
		Expect(os.RemoveAll(filepath.Join(root, "services"))).To(Succeed())
		writeFixture(filepath.Join(root, "sample.go"), "package sample\nfunc Answer() int {return 42}")
		args := []string{"-backend", root, "-core", root, "-generator-script", generator, "-reports", reports, "-package", "core"}
		Expect(coverage.Run(args)).To(MatchError("one or more selected unit coverage runs failed"))
		Expect(coverage.Run(append(args, "-mode", "--report-only"))).To(Succeed())
		writeFixture(filepath.Join(root, "sample.go"), "package sample\nconst Answer = 42")
		Expect(coverage.Run(append(args, "-mode", "--report-only"))).To(MatchError("one or more selected unit coverage runs failed"))
	})
})

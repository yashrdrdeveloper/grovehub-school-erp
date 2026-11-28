package coverage_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/testing/coverage"
)

var _ = Describe("Executable coverage", func() {
	It("covers the real command exit path and includes main in the strict denominator", func() {
		root, generator, reports, _ := coverageFixture()
		writeFixture(filepath.Join(root, "cmd", "main.go"), "package main\nimport \"os\"\nfunc main(){os.Exit(2)}\n")
		args := []string{"-backend", root, "-core", root, "-generator-script", generator, "-reports", reports, "-package", "cmd"}
		Expect(coverage.Run(args)).To(Succeed())
		profile := filepath.Join(root, ".generated", "unit-coverage", "cmd.out")
		total, covered, err := coverage.Score(profile)
		Expect(err).NotTo(HaveOccurred())
		Expect(total).To(Equal(int64(1)))
		Expect(covered).To(Equal(total))

		// A zero exit also works when explicitly expected by a caller.
		writeFixture(filepath.Join(root, "cmd", "main.go"), "package main\nimport \"os\"\nfunc main(){os.Exit(0)}\n")
		writeFixture(profile, "mode: atomic\n")
		Expect(coverage.CoverEntrypoint(root, profile, "./cmd", nil, 0)).To(Succeed())
		Expect(coverage.CoverEntrypoint(root, profile, "./cmd", nil, 2)).To(MatchError(ContainSubstring("expected 2")))
		writeFixture(filepath.Join(root, "cmd", "main.go"), "package main\nfunc invalid(")
		Expect(coverage.CoverEntrypoint(root, profile, "./cmd", nil, 2)).NotTo(Succeed())
		Expect(coverage.CoverEntrypoint(root, filepath.Join(root, "missing", "profile"), "./cmd", nil, 2)).NotTo(Succeed())
	})

	It("propagates missing executable and conversion failures", func() {
		root, _, _, _ := coverageFixture()
		profile := filepath.Join(root, "profile")
		writeFixture(profile, "mode: atomic\n")
		commandFixture(root, "exit 0\n")
		Expect(coverage.CoverEntrypoint(root, profile, "./cmd", nil, 2)).NotTo(Succeed())

		commandFixture(root, "case \"$1\" in\nbuild) while [ \"$#\" -gt 0 ]; do if [ \"$1\" = -o ]; then shift; printf '#!/bin/sh\\nexit 2\\n' > \"$1\"; chmod +x \"$1\"; exit 0; fi; shift; done; exit 1 ;;\ntool) exit 17 ;;\nesac\n")
		Expect(coverage.CoverEntrypoint(root, profile, "./cmd", nil, 2)).To(MatchError(ContainSubstring("exit status 17")))
	})

	It("validates both profiles without silently losing uncovered statements", func() {
		root := GinkgoT().TempDir()
		unit := filepath.Join(root, "unit.out")
		entrypoint := filepath.Join(root, "entrypoint.out")
		Expect(coverage.MergeEntrypointProfile(unit, entrypoint)).NotTo(Succeed())
		writeFixture(unit, "mode: atomic\nfixture.go:1.1,1.2 1 0\nfixture.go:2.1,2.2 1 0\n")
		Expect(coverage.MergeEntrypointProfile(unit, entrypoint)).NotTo(Succeed())
		for _, invalid := range []string{"mode: count\nfixture.go:1.1,1.2 1 1\n", "mode: atomic\n"} {
			writeFixture(entrypoint, invalid)
			Expect(coverage.MergeEntrypointProfile(unit, entrypoint)).NotTo(Succeed())
		}
		writeFixture(entrypoint, "mode: atomic\nfixture.go:1.1,1.2 1 1\n")
		Expect(coverage.MergeEntrypointProfile(unit, entrypoint)).To(Succeed())
		total, covered, err := coverage.Score(unit)
		Expect(err).NotTo(HaveOccurred())
		Expect(total).To(Equal(int64(2)))
		Expect(covered).To(Equal(int64(1)))

		saved, err := os.ReadFile(unit)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(saved)).To(ContainSubstring("fixture.go:2.1,2.2 1 0"))
	})
})

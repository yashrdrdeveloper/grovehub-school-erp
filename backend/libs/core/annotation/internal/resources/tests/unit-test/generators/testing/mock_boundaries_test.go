package testgen_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	testgen "grovehub/libs/core/annotation/internal/src/generators/testing"
)

var _ = Describe("Test generator subprocess and output boundaries", func() {
	It("invokes mockgen from the owning module and propagates child errors", func() {
		root := GinkgoT().TempDir()
		writeTestSource(filepath.Join(root, "go.mod"), "module example.com/project\n")
		source := filepath.Join(root, "tests", "mocks", "store.go")
		writeTestSource(source, "package mocks\n// @Mock\ntype Store interface{ Read() string }")
		bin := filepath.Join(root, "bin", "go")
		log := filepath.Join(root, "mockgen.log")
		GinkgoT().Setenv("MOCKGEN_FIXTURE_LOG", log)
		writeTestSource(bin, "#!/bin/sh\nset -eu\nprintf '%s\\n' \"$PWD\" \"$@\" > \"$MOCKGEN_FIXTURE_LOG\"\nfor arg in \"$@\"; do case \"$arg\" in -destination=*) printf 'package mocks\\n' > \"${arg#-destination=}\";; esac; done\n")
		Expect(os.Chmod(bin, 0755)).To(Succeed())
		GinkgoT().Setenv("PATH", filepath.Dir(bin)+string(os.PathListSeparator)+os.Getenv("PATH"))
		Expect(testgen.GenerateAll(root, root)).To(Succeed())
		output := filepath.Join(root, ".generated", "tests", "mocks", "zz_mocks.gen.go")
		Expect(output).To(BeAnExistingFile())
		args, err := os.ReadFile(log)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(args)).To(ContainSubstring("run\ngo.uber.org/mock/mockgen\n-source=" + source + "\n-destination=" + output + "\n-package=mocks"))
		writeTestSource(bin, "#!/bin/sh\necho 'mockgen fixture failure' >&2\nexit 17\n")
		Expect(testgen.GenerateAll(root, root)).To(MatchError(ContainSubstring("mockgen fixture failure")))
		Expect(os.Remove(source)).To(Succeed())
		Expect(testgen.GenerateAll(root, root)).To(Succeed())
		Expect(output).NotTo(BeAnExistingFile())
	})
	It("reports blocked output directories and invalid stale outputs", func() {
		root := GinkgoT().TempDir()
		dir := filepath.Join(root, "tests", "unit-test", "sample")
		source := filepath.Join(dir, "sample_test.go")
		writeTestSource(source, "package sample_test\n// @Testing suite=Sample\nfunc specs(){}")
		out := filepath.Join(root, ".generated", "tests", "unit-test", "sample", "zz_suite_test.go")
		Expect(os.MkdirAll(out, 0755)).To(Succeed())
		Expect(testgen.GenerateAll(root, root)).NotTo(Succeed())
		Expect(os.Remove(out)).To(Succeed())
		Expect(os.Remove(filepath.Dir(out))).To(Succeed())
		writeTestSource(filepath.Dir(out), "blocked")
		Expect(testgen.GenerateAll(root, root)).NotTo(Succeed())
		Expect(os.Remove(filepath.Dir(out))).To(Succeed())
		Expect(os.Remove(source)).To(Succeed())
		Expect(os.MkdirAll(filepath.Join(out, "child"), 0755)).To(Succeed())
		Expect(testgen.GenerateAll(root, root)).NotTo(Succeed())
	})
	It("reports a blocked mock output and duplicate directives on mock interfaces", func() {
		root := GinkgoT().TempDir()
		source := filepath.Join(root, "tests", "mocks", "store.go")
		writeTestSource(source, "package mocks\n// @Mock\n// @Mock\ntype Store interface{}")
		Expect(testgen.GenerateAll(root, root)).To(MatchError(ContainSubstring("duplicate @Mock")))
		writeTestSource(source, "package mocks\n// @Mock\ntype Store interface{}")
		writeTestSource(filepath.Join(root, ".generated", "tests", "mocks"), "blocked")
		Expect(testgen.GenerateAll(root, root)).NotTo(Succeed())
		Expect(os.Remove(filepath.Join(root, ".generated", "tests", "mocks"))).To(Succeed())
		Expect(os.Remove(source)).To(Succeed())
		Expect(os.MkdirAll(filepath.Join(root, ".generated", "tests", "mocks", "zz_mocks.gen.go", "child"), 0755)).To(Succeed())
		Expect(testgen.GenerateAll(root, root)).NotTo(Succeed())
	})
	It("handles literal metacharacter paths and reports filesystem errors", func() {
		root := GinkgoT().TempDir()
		Expect(os.MkdirAll(filepath.Join(root, "tests", "unit-test", "["), 0755)).To(Succeed())
		Expect(testgen.GenerateAll(root, root)).To(Succeed())
		Expect(os.RemoveAll(filepath.Join(root, "tests"))).To(Succeed())
		Expect(os.Symlink("tests", filepath.Join(root, "tests"))).To(Succeed())
		Expect(testgen.GenerateAll(root, root)).NotTo(Succeed())
	})
})

package cli_test

import (
	"bytes"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/cli"
)

var _ = Describe("Run", func() {
	It("lists all fourteen processors in deterministic order", func() {
		var stdout, stderr bytes.Buffer

		Expect(cli.Run([]string{"--help"}, &stdout, &stderr)).To(Equal(0))
		Expect(stderr.String()).To(BeEmpty())

		for _, name := range []string{
			"check", "clients", "coverage", "entities", "events", "expand",
			"gormtable", "inventory", "mappers", "overlay", "patchdto", "routes",
			"search", "testing",
		} {
			Expect(stdout.String()).To(ContainSubstring("  " + name))
		}
	})

	It("rejects unknown commands and incorrect argument counts", func() {
		var stdout, stderr bytes.Buffer

		Expect(cli.Run(nil, &stdout, &stderr)).To(Equal(2))
		Expect(cli.Run([]string{"missing"}, &stdout, &stderr)).To(Equal(2))
		Expect(cli.Run([]string{"routes"}, &stdout, &stderr)).To(Equal(2))
		Expect(cli.Run([]string{"overlay", "one"}, &stdout, &stderr)).To(Equal(2))
		Expect(cli.Run([]string{"routes", "--help"}, &stdout, &stderr)).To(Equal(0))
	})

	It("dispatches the inventory processor and writes JSON", func() {
		var stdout, stderr bytes.Buffer
		root := GinkgoT().TempDir()

		Expect(os.WriteFile(filepath.Join(root, "demo.go"), []byte("package demo\nfunc Run() {}\n"), 0600)).To(Succeed())

		Expect(cli.Run([]string{"inventory", root}, &stdout, &stderr)).To(Equal(0))
		Expect(stdout.String()).To(ContainSubstring("\"sourceFiles\": 1"))
		Expect(cli.Run([]string{"inventory", root, "--strict"}, &stdout, &stderr)).To(Equal(1))
	})

	It("dispatches testing and overlay with explicit roots", func() {
		var stdout, stderr bytes.Buffer
		root := GinkgoT().TempDir()

		Expect(cli.Run([]string{"testing", root, root}, &stdout, &stderr)).To(Equal(0))
		Expect(cli.Run([]string{"overlay", root, root}, &stdout, &stderr)).To(Equal(0))
		Expect(filepath.Join(root, ".generated", "overlay.json")).To(BeAnExistingFile())
	})

	It("accepts coverage flags after the command name", func() {
		var stdout, stderr bytes.Buffer

		Expect(cli.Run([]string{"coverage", "--help"}, &stdout, &stderr)).To(Equal(0))
		Expect(cli.Run([]string{"coverage", "-mode", "--report-only"}, &stdout, &stderr)).To(Equal(1))
		Expect(stderr.String()).To(ContainSubstring("core, generator-script and reports are required"))
	})
})

package pack_test

import (
	"archive/zip"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/tools/modulepack/internal/src/operations/pack"
)

const (
	modulePath = "example.com/School/Core"
	version    = "v0.1.0-ci.123"
	stamp      = "2026-11-28T10:30:00+05:30"
	moduleFile = "module " + modulePath + "\n\ngo 1.27.0\n"
)

func fixture() (string, string, []string) {
	source := GinkgoT().TempDir()
	proxy := GinkgoT().TempDir()
	Expect(os.WriteFile(filepath.Join(source, "go.mod"), []byte(moduleFile), 0644)).To(Succeed())
	Expect(os.WriteFile(filepath.Join(source, "core.go"), []byte("package core\n"), 0644)).To(Succeed())
	return source, proxy, []string{source, proxy, version, stamp}
}

func versionDirectory(proxy string) string {
	return filepath.Join(proxy, "example.com", "!school", "!core", "@v")
}

var _ = Describe("Run", func() {
	It("creates an escaped Go proxy path with valid archive, module, info and list artifacts", func() {
		_, proxy, args := fixture()
		Expect(pack.Run(args)).To(Succeed())
		dir := versionDirectory(proxy)
		mod, err := os.ReadFile(filepath.Join(dir, version+".mod"))
		Expect(err).NotTo(HaveOccurred())
		Expect(string(mod)).To(Equal(moduleFile))
		list, err := os.ReadFile(filepath.Join(dir, "list"))
		Expect(err).NotTo(HaveOccurred())
		Expect(string(list)).To(Equal(version + "\n"))
		info, err := os.ReadFile(filepath.Join(dir, version+".info"))
		Expect(err).NotTo(HaveOccurred())
		var metadata struct {
			Version string
			Time    time.Time
		}
		Expect(json.Unmarshal(info, &metadata)).To(Succeed())
		Expect(metadata.Version).To(Equal(version))
		Expect(metadata.Time.Format(time.RFC3339)).To(Equal("2026-11-28T05:00:00Z"))
		archive, err := zip.OpenReader(filepath.Join(dir, version+".zip"))
		Expect(err).NotTo(HaveOccurred())
		defer archive.Close()
		names := []string{}
		for _, file := range archive.File {
			names = append(names, file.Name)
		}
		Expect(names).To(ContainElements(modulePath+"@"+version+"/go.mod", modulePath+"@"+version+"/core.go"))
	})

	It("refuses to overwrite an existing release archive", func() {
		_, proxy, args := fixture()
		Expect(pack.Run(args)).To(Succeed())
		path := filepath.Join(versionDirectory(proxy), version+".zip")
		before, err := os.ReadFile(path)
		Expect(err).NotTo(HaveOccurred())
		Expect(pack.Run(args)).NotTo(Succeed())
		after, err := os.ReadFile(path)
		Expect(err).NotTo(HaveOccurred())
		Expect(after).To(Equal(before))
	})

	It("rejects incorrect invocation, missing or malformed module files, and replacements", func() {
		Expect(pack.Run(nil)).To(MatchError("usage: modulepack SOURCE PROXY VERSION RFC3339_TIME"))
		source, _, args := fixture()
		mod := filepath.Join(source, "go.mod")
		Expect(os.Remove(mod)).To(Succeed())
		Expect(pack.Run(args)).NotTo(Succeed())
		Expect(os.WriteFile(mod, []byte("not a module directive\n"), 0644)).To(Succeed())
		Expect(pack.Run(args)).NotTo(Succeed())
		for _, content := range []string{"go 1.27.0\n", moduleFile + "replace example.com/dependency => ../dependency\n"} {
			Expect(os.WriteFile(mod, []byte(content), 0644)).To(Succeed())
			Expect(pack.Run(args)).To(MatchError("module must declare its path and have no replacements"))
		}
	})

	It("rejects invalid versions and timestamps before publishing", func() {
		_, proxy, args := fixture()
		args[2] = "not-a-version"
		Expect(pack.Run(args)).NotTo(Succeed())
		args[2] = version
		args[3] = "not-a-timestamp"
		Expect(pack.Run(args)).NotTo(Succeed())
		Expect(versionDirectory(proxy)).NotTo(BeADirectory())
	})

	It("propagates directory and artifact creation failures", func() {
		_, proxy, args := fixture()
		blocked := filepath.Join(proxy, "blocked")
		Expect(os.WriteFile(blocked, []byte("file"), 0644)).To(Succeed())
		args[1] = blocked
		Expect(pack.Run(args)).NotTo(Succeed())
		for _, name := range []string{version + ".zip", version + ".mod", version + ".info", "list"} {
			_, proxy, args := fixture()
			dir := versionDirectory(proxy)
			Expect(os.MkdirAll(filepath.Join(dir, name), 0755)).To(Succeed())
			Expect(pack.Run(args)).NotTo(Succeed())
		}
	})
})

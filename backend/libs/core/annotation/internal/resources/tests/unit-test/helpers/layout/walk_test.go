package layout_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/helpers/layout"
)

var _ = Describe("WalkOwned", func() {
	It("returns paths relative to their owner and skips directory symlinks and excluded directories", func() {
		root := GinkgoT().TempDir()
		for _, rel := range []string{"source/record", "source/vendor", "source/.generated"} {
			Expect(os.MkdirAll(filepath.Join(root, filepath.FromSlash(rel)), 0755)).To(Succeed())
		}
		Expect(os.WriteFile(filepath.Join(root, "source", "record", "source.go"), []byte("package record"), 0644)).To(Succeed())
		Expect(os.Symlink("record", filepath.Join(root, "source", "alias"))).To(Succeed())
		var dirs []string
		Expect(layout.WalkOwnedDirectories(root, "source", map[string]bool{"vendor": true, ".generated": true}, func(path, relative string) error {
			Expect(path).To(Equal(filepath.Join(root, relative)))
			dirs = append(dirs, filepath.ToSlash(relative))
			return nil
		})).To(Succeed())
		Expect(dirs).To(Equal([]string{"source", "source/record"}))
		var files []string
		Expect(layout.WalkOwned(root, "source", func(_ string, relative string, entry fs.DirEntry) error {
			if !entry.IsDir() {
				files = append(files, filepath.ToSlash(relative))
			}
			return nil
		})).To(Succeed())
		Expect(files).To(ConsistOf("source/alias", "source/record/source.go"))
		Expect(layout.WalkOwned(root, "source/alias", func(string, string, fs.DirEntry) error { Fail("must not follow root symlink"); return nil })).To(Succeed())
	})
	It("propagates input, visitor and walk errors", func() {
		root := GinkgoT().TempDir()
		visit := func(string, string, fs.DirEntry) error { Fail("must not visit"); return nil }
		Expect(layout.WalkOwned(root, "missing", visit)).To(Succeed())
		Expect(layout.WalkOwned(root, "../outside", visit)).NotTo(Succeed())
		Expect(os.Symlink("loop", filepath.Join(root, "loop"))).To(Succeed())
		Expect(layout.WalkOwned(root, "loop", visit)).NotTo(Succeed())
		Expect(layout.WalkOwned(root, "loop/child", visit)).NotTo(Succeed())
		Expect(os.Symlink("missing", filepath.Join(root, "dangling"))).To(Succeed())
		Expect(layout.WalkOwned(root, "dangling", visit)).To(Succeed())
		Expect(os.WriteFile(filepath.Join(root, "file"), []byte("plain"), 0644)).To(Succeed())
		Expect(layout.WalkOwned(root, "file", visit)).To(Succeed())
		Expect(os.Mkdir(filepath.Join(root, "source"), 0755)).To(Succeed())
		failure := errors.New("visit failed")
		Expect(layout.WalkOwnedDirectories(root, "source", nil, func(string, string) error { return failure })).To(MatchError(failure))
		if os.Geteuid() != 0 {
			dir := filepath.Join(root, "source", "unreadable")
			Expect(os.Mkdir(dir, 0000)).To(Succeed())
			defer os.Chmod(dir, 0755)
			Expect(layout.WalkOwned(root, "source", func(string, string, fs.DirEntry) error { return nil })).NotTo(Succeed())
		}
	})
})

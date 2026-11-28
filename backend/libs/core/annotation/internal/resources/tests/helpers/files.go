// Package helpers provides temporary source fixtures for annotation unit tests.
package helpers

import (
	"os"
	"path/filepath"

	. "github.com/onsi/gomega"
)

func Write(root, relative, source string) string {
	path := filepath.Join(root, filepath.FromSlash(relative))
	Expect(os.MkdirAll(filepath.Dir(path), 0700)).To(Succeed())
	Expect(os.WriteFile(path, []byte(source), 0600)).To(Succeed())

	return path
}

func Read(path string) string {
	data, err := os.ReadFile(path)
	Expect(err).NotTo(HaveOccurred())

	return string(data)
}

package clients_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	filefixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
	"grovehub/libs/core/annotation/internal/src/generators/clients"
)

var _ = Describe("clients filesystem discovery", func() {
	It("rejects absent service roots and handles empty or unrelated services", func() {
		root := GinkgoT().TempDir()

		Expect(clients.GenerateAll(root)).NotTo(Succeed())
		filefixture.Write(root, "services/README.md", "not a service")
		Expect(os.MkdirAll(filepath.Join(root, "services", "empty"), 0755)).To(Succeed())
		Expect(clients.GenerateAll(root)).To(Succeed())
	})
	It("reports invalid maintained source rather than generating output", func() {
		root := GinkgoT().TempDir()

		filefixture.Write(root, "services/demo/internal/app/client/bad.go", "package broken\nfunc malformed(")
		filefixture.Write(root, "services/demo/client/bad.go", "package broken\nfunc malformed(")
		Expect(clients.GenerateAll(root)).NotTo(Succeed())
	})
})

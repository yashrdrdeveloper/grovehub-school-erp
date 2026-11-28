package routes_test

import (
    "os"
    "path/filepath"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/routes"
)

var _ = Describe("GenerateAll", func() {
    It("maps protected controllers to restricted project routes and rejects duplicate routes", func() {
        t := GinkgoT()
        backend, output := testBackend(t)
        normal := generatedRoutes(t, backend, output)
        Expect(normal).To(ContainSubstring("GET /foundation/protected/persons/{gid}"))
        Expect(os.WriteFile(filepath.Join(backend, ".restricted-routes"), nil, 0600)).To(Succeed())
        restricted := generatedRoutes(t, backend, output)
        Expect(restricted).To(ContainSubstring("GET /foundation/restricted/persons/{gid}"))
        Expect(restricted).To(ContainSubstring("func registerProtectedPersonRoutes("))
        Expect(restricted).NotTo(ContainSubstring("GET /foundation/protected/persons/{gid}"))
        owner := filepath.Join(backend, "services", "foundation")
        writeTestFile(t, filepath.Join(owner, "controller", "public", "duplicate.go"), "package public"+personController[len("package protected"):])
        Expect(routes.GenerateAll(backend)).NotTo(Succeed())
    })
})

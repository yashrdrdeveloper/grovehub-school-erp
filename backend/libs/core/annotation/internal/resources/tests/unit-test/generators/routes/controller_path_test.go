package routes_test

import (
"os"
"path/filepath"
"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/routes"
. "github.com/onsi/ginkgo/v2"
. "github.com/onsi/gomega"
)

var _=Describe("Controller path errors",func(){It("reports a controller symlink loop",func(){root:=GinkgoT().TempDir();Expect(os.WriteFile(filepath.Join(root,"go.mod"),[]byte("module grovehub\n"),0644)).To(Succeed());dir:=filepath.Join(root,"services","demo");Expect(os.MkdirAll(dir,0755)).To(Succeed());Expect(os.Symlink("controller",filepath.Join(dir,"controller"))).To(Succeed());Expect(routes.GenerateAll(root)).NotTo(Succeed())})})

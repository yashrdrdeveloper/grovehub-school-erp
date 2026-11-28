package inventory_test
import (
 "os"
 "path/filepath"
 . "github.com/onsi/ginkgo/v2"
 . "github.com/onsi/gomega"
 "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/testing/inventory"
)
var _ = Describe("Launcher inventory",func(){
 It("maps internal/src to internal/resources/tests",func(){
  owner:=filepath.Join(GinkgoT().TempDir(),"launcher")
  source:=filepath.Join(owner,"internal","src","config","config.go")
  mirror:=filepath.Join(owner,"internal","resources","tests","unit-test","config","config_test.go")
  for _,p:=range []string{source,mirror}{Expect(os.MkdirAll(filepath.Dir(p),0755)).To(Succeed())}
  Expect(os.WriteFile(source,[]byte("package config\nfunc Run(){}\n"),0600)).To(Succeed())
  Expect(os.WriteFile(mirror,[]byte("package config_test\n"),0600)).To(Succeed())
  report,err:=inventory.Inspect(owner)
  Expect(err).NotTo(HaveOccurred())
  Expect(report.Functions).To(Equal(1))
  Expect(report.MirroredFunctions).To(Equal(1))
 })
})

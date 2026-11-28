package testgen_test
import (
 "os"
 "path/filepath"
 . "github.com/onsi/ginkgo/v2"
 . "github.com/onsi/gomega"
 "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/testing"
)
var _ = Describe("Launcher test generation",func(){
 It("generates @Testing suites in the Launcher owner",func(){
  root:=GinkgoT().TempDir()
  app,core:=filepath.Join(root,"app"),filepath.Join(root,"core")
  owner:=filepath.Join(app,"launcher")
  rel:=filepath.Join("internal","resources","tests","unit-test","controller")
  dir:=filepath.Join(owner,rel)
  for _,p:=range []string{dir,core,filepath.Join(app,"services")}{Expect(os.MkdirAll(p,0755)).To(Succeed())}
  Expect(os.WriteFile(filepath.Join(dir,"health_test.go"),[]byte("package controller_test\n// @Testing suite=Health\nfunc specs(){}\n"),0600)).To(Succeed())
  Expect(testgen.GenerateAll(app,core)).To(Succeed())
  data,err:=os.ReadFile(filepath.Join(owner,".generated",rel,"zz_suite_test.go"))
  Expect(err).NotTo(HaveOccurred())
  Expect(string(data)).To(ContainSubstring("func TestHealth"))
 })
})

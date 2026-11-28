package testgen_test

import (
"os"
"path/filepath"
"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/testing"
. "github.com/onsi/ginkgo/v2"
. "github.com/onsi/gomega"
)

var _=Describe("Unannotated test declarations",func(){It("ignores types and functions without testing annotations",func(){root:=GinkgoT().TempDir();dir:=filepath.Join(root,"tests","unit-test","plain");Expect(os.MkdirAll(dir,0755)).To(Succeed());Expect(os.WriteFile(filepath.Join(dir,"plain_test.go"),[]byte("package plain_test\ntype Plain struct{}\nfunc helper() {}"),0644)).To(Succeed());Expect(testgen.GenerateAll(root,root)).To(Succeed())})})

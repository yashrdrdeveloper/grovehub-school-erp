package inventory_test

import (
"os"
"path/filepath"
"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/testing/inventory"
. "github.com/onsi/ginkgo/v2"
. "github.com/onsi/gomega"
)

var _=Describe("Inventory filesystem failures",func(){It("reports missing owners",func(){_,err:=inventory.Inspect(filepath.Join(GinkgoT().TempDir(),"missing"));Expect(err).To(HaveOccurred())});It("reports malformed maintained source",func(){root:=GinkgoT().TempDir();Expect(os.WriteFile(filepath.Join(root,"bad.go"),[]byte("package"),0644)).To(Succeed());_,err:=inventory.Inspect(root);Expect(err).To(HaveOccurred())});It("reports unreadable test trees",func(){if os.Geteuid()==0{Skip("root bypasses permission checks")};root:=GinkgoT().TempDir();dir:=filepath.Join(root,"tests","unit-test");Expect(os.MkdirAll(dir,0755)).To(Succeed());Expect(os.Chmod(dir,0000)).To(Succeed());DeferCleanup(func(){_ = os.Chmod(dir,0755)});_,err:=inventory.Inspect(root);Expect(err).To(HaveOccurred())})})

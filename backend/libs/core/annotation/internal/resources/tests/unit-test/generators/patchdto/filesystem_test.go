package patchdto_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/patchdto"
	filefixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("patchdto filesystem discovery",func(){
	It("rejects absent service roots and handles empty or unrelated services",func(){
		root:=GinkgoT().TempDir()
		
		Expect(patchdto.GenerateAll(root)).NotTo(Succeed())
		filefixture.Write(root,"services/README.md","not a service")
		Expect(os.MkdirAll(filepath.Join(root,"services","empty"),0755)).To(Succeed())
		Expect(patchdto.GenerateAll(root)).To(Succeed())
	})
	It("reports invalid maintained source rather than generating output",func(){
		root:=GinkgoT().TempDir()
		
		filefixture.Write(root,"services/demo/internal/app/dto/bad.go","package broken\nfunc malformed(")
		
		Expect(patchdto.GenerateAll(root)).NotTo(Succeed())
	})
})

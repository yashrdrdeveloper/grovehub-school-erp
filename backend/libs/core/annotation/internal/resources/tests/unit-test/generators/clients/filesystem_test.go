package clients_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/clients"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("clients filesystem discovery",func(){
	It("rejects absent service roots and handles empty or unrelated services",func(){
		root:=GinkgoT().TempDir()
		
		Expect(clients.GenerateAll(root)).NotTo(Succeed())
		fixture.Write(root,"services/README.md","not a service")
		Expect(os.MkdirAll(filepath.Join(root,"services","empty"),0755)).To(Succeed())
		Expect(clients.GenerateAll(root)).To(Succeed())
	})
	It("reports invalid maintained source rather than generating output",func(){
		root:=GinkgoT().TempDir()
		
		fixture.Write(root,"services/demo/internal/app/client/bad.go","package broken\nfunc malformed(")
		fixture.Write(root,"services/demo/client/bad.go","package broken\nfunc malformed(")
		Expect(clients.GenerateAll(root)).NotTo(Succeed())
	})
})

package mappers_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/mappers"
	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("mappers filesystem discovery",func(){
	It("rejects absent service roots and handles empty or unrelated services",func(){
		root:=GinkgoT().TempDir()
		Expect(mappers.GenerateAll(root)).NotTo(Succeed())
		fixture.Write(root,"go.mod","module example.com/project\n")
		Expect(mappers.GenerateAll(root)).NotTo(Succeed())
		fixture.Write(root,"services/README.md","not a service")
		Expect(os.MkdirAll(filepath.Join(root,"services","empty"),0755)).To(Succeed())
		Expect(mappers.GenerateAll(root)).To(Succeed())
	})
	It("reports invalid maintained source rather than generating output",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"go.mod","module example.com/project\n")
		fixture.Write(root,"services/demo/internal/app/mappers/bad.go","package broken\nfunc malformed(")
		
		Expect(mappers.GenerateAll(root)).NotTo(Succeed())
	})
})

package entities_test

import (
	"path/filepath"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/entities"
	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
)

var _ = Describe("GenerateAll",func(){
	It("generates only annotated non-entity accessors under the service-owned output directory",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"services/demo/internal/app/entity/record.go","package entity\n// @Getter\ntype Record struct{ Name string }\n")
		fixture.Write(root,"services/demo/internal/app/entity/persisted.go","package entity\n// @Entity demo.persisted\ntype Persisted struct{}\n")
		Expect(entities.GenerateAll(root)).To(Succeed())
		output:=filepath.Join(root,"services/demo/.generated/internal/app/entity")
		Expect(fixture.Read(filepath.Join(output,"record_gombok.go"))).To(ContainSubstring("GetName() string"))
		Expect(filepath.Join(output,"persisted_gombok.go")).NotTo(BeAnExistingFile())
	})
})

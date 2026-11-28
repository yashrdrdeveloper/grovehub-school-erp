package mappers_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/mappers"
	sourceFixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("Internal mapper layout",func(){
	It("writes resource-local output and removes or rejects stale generated files",func(){
		root:=GinkgoT().TempDir()
		sourceFixture.Write(root,"go.mod","module example.com/project\n")
		owner:=filepath.Join(root,"services","demo")
		sourceFixture.Write(owner,"internal/app/dto/record.go","package dto\ntype Record struct{ Name string }")
		sourceFixture.Write(owner,"internal/app/mappers/record/source.go","package record\nimport dto \"example.com/project/services/demo/internal/app/dto\"\n// @Mapper impl\ntype Mapper interface{\n// @Mapping Name=input\nMap(input string) dto.Record\n}")
		Expect(mappers.GenerateAll(root)).To(Succeed())
		out:=filepath.Join(owner,".generated","internal","app","mappers","record","impl","zz_mappers.gen.go")
		Expect(sourceFixture.Read(out)).To(ContainSubstring("Name: input"))
		Expect(os.Remove(out)).To(Succeed())
		Expect(os.MkdirAll(filepath.Join(out,"child"),0755)).To(Succeed())
		Expect(mappers.GenerateAll(root)).NotTo(Succeed())
		sourceFixture.Write(owner,"internal/app/mappers/record/source.go","package record\ntype Mapper interface{ Maintained() }")
		Expect(mappers.GenerateAll(root)).NotTo(Succeed())
		Expect(os.RemoveAll(out)).To(Succeed())
		Expect(mappers.GenerateAll(root)).To(Succeed())
		Expect(out).NotTo(BeAnExistingFile())
	})
})

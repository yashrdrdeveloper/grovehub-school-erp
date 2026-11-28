package expand_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/expand"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("Entity annotation expansion",func(){
	It("preserves nonentity source and explicit setters while expanding remaining accessors",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"ignored_test.go","invalid source")
		fixture.Write(root,"ignored_gombok.go","invalid source")
		plain:=fixture.Write(root,"plain.go","package entity\ntype Plain string\nconst x=1\nfunc Helper(){}\n")
		path:=fixture.Write(root,"record.go","package entity\n// description\n// @Entity demo.records\ntype Record struct{ Name string }\nfunc (r Record) SetName(name string){}\n")
		Expect(expand.Expand(root)).To(Succeed())
		data:=fixture.Read(path)
		Expect(data).To(ContainSubstring("@Builder"))
		Expect(data).To(ContainSubstring("@Getter"))
		Expect(data).NotTo(ContainSubstring("@Setter"))
		Expect(fixture.Read(plain)).To(ContainSubstring("type Plain string"))
		Expect(expand.Expand(root)).To(Succeed())
		Expect(fixture.Read(path)).To(Equal(data))
		fixture.Write(root,"broken.go","package entity\nfunc broken(")
		Expect(expand.Expand(root)).NotTo(Succeed())
	})
})

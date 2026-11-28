package expand_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/expand"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _ = Describe("Expand",func(){
	It("preserves manually implemented accessor categories when expanding staging annotations",func(){
		root:=GinkgoT().TempDir()
		path:=fixture.Write(root,"record.go","package entity\n// @Entity demo.records\ntype Record struct{ Name string }\nfunc (r *Record) GetName() string{return r.Name}\n")
		Expect(expand.Expand(root)).To(Succeed())
		content:=fixture.Read(path)
		Expect(content).To(ContainSubstring("// @Builder"))
		Expect(content).To(ContainSubstring("// @Setter"))
		Expect(content).NotTo(ContainSubstring("// @Getter"))
		Expect(content).NotTo(ContainSubstring("// @Entity"))
	})
	It("rejects entity annotations with multiple table arguments",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"record.go","package entity\n// @Entity demo.records extra\ntype Record struct{}\n")
		Expect(expand.Expand(root)).To(MatchError(ContainSubstring("single schema.table name")))
	})
})

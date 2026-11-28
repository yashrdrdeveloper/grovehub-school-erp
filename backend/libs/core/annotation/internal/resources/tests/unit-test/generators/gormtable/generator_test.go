package gormtable_test

import (
	"path/filepath"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/gormtable"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _ = Describe("GenerateInto",func(){
	It("emits deterministic table methods and removes stale generated output",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"record.go","package entity\n// @Entity demo.records mode=table\ntype Record struct{}\n")
		out:=filepath.Join(GinkgoT().TempDir(),"zz_tables.gen.go")
		Expect(gormtable.GenerateInto(root,out)).To(Succeed())
		content:=fixture.Read(out)
		Expect(content).To(ContainSubstring(`func (Record) TableName() string { return "demo.records" }`))
		Expect(gormtable.GenerateInto(root,out)).To(Succeed())
		Expect(fixture.Read(out)).To(Equal(content))
		fixture.Write(root,"record.go","package entity\ntype Record struct{}\n")
		Expect(gormtable.GenerateInto(root,out)).To(Succeed())
		Expect(out).NotTo(BeAnExistingFile())
	})
	It("rejects handwritten table names on annotated types",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"record.go","package entity\n// @Entity demo.records\ntype Record struct{}\nfunc (Record) TableName()string{return \"other\"}\n")
		Expect(gormtable.GenerateInto(root,filepath.Join(root,"out.go"))).To(MatchError(ContainSubstring("declares TableName manually")))
	})
})

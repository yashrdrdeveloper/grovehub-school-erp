package gormtable_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/gormtable"
	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("Table mapping source validation",func(){
	It("rejects missing source, invalid declarations, conflicting annotations and mixed packages",func(){
		root:=GinkgoT().TempDir()
		out:=filepath.Join(GinkgoT().TempDir(),"zz_tables.gen.go")
		Expect(gormtable.GenerateInto(filepath.Join(root,"missing"),out)).NotTo(Succeed())
		Expect(gormtable.GenerateInto(root,out)).To(MatchError(ContainSubstring("no Go source files")))
		for _,test:=range []struct{source,reason string}{
			{"// @Entity demo.records extra\ntype Record struct{}","invalid @Entity options"},
			{"// @Entity demo.Bad\ntype Record struct{}","invalid @Entity"},
			{"// @Entity demo.one\n// @Entity demo.two\ntype Record struct{}","conflicting table declarations"},
			{"func broken(","expected"},
		}{
			fixture.Write(root,"record.go","package entity\n"+test.source)
			Expect(gormtable.GenerateInto(root,out)).To(MatchError(ContainSubstring(test.reason)))
		}
		fixture.Write(root,"record.go","package entity")
		fixture.Write(root,"other.go","package other")
		Expect(gormtable.GenerateInto(root,out)).To(MatchError(ContainSubstring("mixed packages")))
	})
	It("ignores unannotated declarations and generated sources and rejects blocked output paths",func(){
		root,out:=GinkgoT().TempDir(),GinkgoT().TempDir()
		fixture.Write(root,"ignored_test.go","bad syntax")
		fixture.Write(root,"ignored.gen.go","bad syntax")
		fixture.Write(root,"ignored_gombok.go","bad syntax")
		fixture.Write(root,"README.md","not Go")
		Expect(os.Mkdir(filepath.Join(root,"directory"),0755)).To(Succeed())
		fixture.Write(root,"record.go","package entity\nconst x=1\ntype Plain string\ntype Unmarked struct{}\n// @Entity demo.records\ntype Record struct{}\nfunc (r *Unmarked) TableName() string{return \"plain\"}\n")
		Expect(gormtable.GenerateInto(root,filepath.Join(out,"zz_tables.gen.go"))).To(Succeed())
		blocked:=fixture.Write(out,"blocked","file")
		Expect(gormtable.GenerateInto(root,filepath.Join(blocked,"out.go"))).NotTo(Succeed())
		Expect(os.Mkdir(filepath.Join(out,"directory.go"),0755)).To(Succeed())
		Expect(gormtable.GenerateInto(root,filepath.Join(out,"directory.go"))).NotTo(Succeed())
	})
})

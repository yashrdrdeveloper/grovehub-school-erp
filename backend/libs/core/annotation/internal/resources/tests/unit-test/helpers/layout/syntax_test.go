package layout_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/helpers/layout"
)

var _=Describe("ParseSource",func(){
	It("decodes aliased imports, ordinary imports and field tags once",func(){
		path:=filepath.Join(GinkgoT().TempDir(),"source.go")
		Expect(os.WriteFile(path,[]byte("package demo\nimport clock \"time\"\nimport \"context\"\ntype Record struct{ Name string `json:\"name\"` }"),0644)).To(Succeed())
		result,err:=layout.ParseSource(path,parser.ParseComments)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.Tree.Name.Name).To(Equal("demo"))
		Expect(result.Imports).To(Equal([]layout.ImportBinding{{Name:"clock",Path:"time",Explicit:true},{Name:"context",Path:"context"}}))
		Expect(result.Tags).To(HaveLen(1))
		for _,tag:=range result.Tags{Expect(tag).To(Equal("json:\"name\""))}
		_,err=layout.ParseSource(filepath.Join(filepath.Dir(path),"missing.go"),0)
		Expect(err).To(HaveOccurred())
	})
	It("rejects invalid caller-created AST literals",func(){
		_,err:=layout.DecodeSource(nil,token.NewFileSet());Expect(err).To(HaveOccurred())
		for _,entry:=range []*ast.ImportSpec{nil,{}, {Path:&ast.BasicLit{Kind:token.STRING,Value:"invalid"}}}{
			_,err:=layout.DecodeSource(&ast.File{Imports:[]*ast.ImportSpec{entry}},token.NewFileSet())
			Expect(err).To(HaveOccurred())
		}
		field:=&ast.Field{Names:[]*ast.Ident{ast.NewIdent("Name")},Type:ast.NewIdent("string"),Tag:&ast.BasicLit{Kind:token.STRING,Value:"invalid"}}
		tree:=&ast.File{Name:ast.NewIdent("demo"),Decls:[]ast.Decl{&ast.GenDecl{Tok:token.TYPE,Specs:[]ast.Spec{&ast.TypeSpec{Name:ast.NewIdent("Record"),Type:&ast.StructType{Fields:&ast.FieldList{List:[]*ast.Field{field,field}}}}}}}}
		_,err=layout.DecodeSource(tree,token.NewFileSet());Expect(err).To(HaveOccurred())
	})
})

var _=Describe("TypeSpecs",func(){
	It("returns only typed declarations",func(){
		item:=&ast.TypeSpec{Name:ast.NewIdent("Record"),Type:ast.NewIdent("string")}
		group:=&ast.GenDecl{Tok:token.TYPE,Specs:[]ast.Spec{&ast.ValueSpec{},item}}
		Expect(layout.TypeSpecs(group)).To(Equal([]*ast.TypeSpec{item}))
		Expect(layout.TypeSpecs(&ast.GenDecl{})).To(BeEmpty())
	})
})

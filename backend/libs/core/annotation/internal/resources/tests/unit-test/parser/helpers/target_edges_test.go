package helpers_test

import (
	"go/ast"
	"go/parser"
	"go/token"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/parser/helpers"
	contract "grovehub/libs/core/annotation/internal/src/contracts"
)

var _=Describe("AST declaration ownership",func(){
	It("distinguishes grouped types, embedded interfaces and file-level declarations",func(){
		src:="package demo\n// group docs\ntype (\n// interface docs\nService interface {\n// embedded docs\nOther\n// method docs\nRead()\n}\n)\n// constants docs\nconst value=1\n"
		file,err:=parser.ParseFile(token.NewFileSet(),"demo.go",src,parser.ParseComments)
		Expect(err).NotTo(HaveOccurred())
		targets:=helpers.CollectTargets(file)
		Expect(targets.For(file.Comments[0])).To(Equal(contract.InterfaceTarget))
		Expect(targets.For(file.Comments[2])).To(Equal(contract.FieldTarget))
		Expect(targets.For(file.Comments[3])).To(Equal(contract.MethodTarget))
		Expect(targets.For(file.Comments[4])).To(Equal(contract.FileTarget))
		group:=&ast.CommentGroup{List:[]*ast.Comment{{Text:"// grouped"}}}
		f:=&ast.File{Name:ast.NewIdent("demo"),Decls:[]ast.Decl{&ast.GenDecl{Doc:group,Tok:token.TYPE,Specs:[]ast.Spec{&ast.TypeSpec{Name:ast.NewIdent("A"),Type:ast.NewIdent("string")},&ast.TypeSpec{Name:ast.NewIdent("B"),Type:ast.NewIdent("int")}}}}}
		Expect(helpers.CollectTargets(f).For(group)).To(Equal(contract.TypeTarget))
	})
})

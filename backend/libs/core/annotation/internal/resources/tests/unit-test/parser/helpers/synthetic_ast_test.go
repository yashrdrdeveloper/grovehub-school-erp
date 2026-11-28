package helpers_test

import (
	"go/ast"
	"go/token"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	contract "grovehub/libs/core/annotation/internal/src/contracts"
	"grovehub/libs/core/annotation/internal/src/parser/helpers"
)

var _ = Describe("Synthetic AST ownership", func() {
	It("handles incomplete caller-created type declarations without panicking", func() {
		group := &ast.CommentGroup{List: []*ast.Comment{{Text: "// partial type"}}}
		f := &ast.File{Name: ast.NewIdent("demo"), Decls: []ast.Decl{&ast.GenDecl{Tok: token.TYPE, Doc: group, Specs: []ast.Spec{&ast.ValueSpec{Names: []*ast.Ident{ast.NewIdent("incomplete")}}}}}}
		Expect(helpers.CollectTargets(f).For(group)).To(Equal(contract.TypeTarget))
	})
})

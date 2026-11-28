package helpers_test

import (
	"go/ast"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/search/helpers"
)

var _ = Describe("ExpressionString", func() {
	It("normalizes type expressions and fails closed for invalid AST nodes", func() {
		Expect(helpers.ExpressionString(&ast.SelectorExpr{X: ast.NewIdent("context"), Sel: ast.NewIdent("Context")})).To(Equal("context.Context"))
		Expect(helpers.ExpressionString(nil)).To(BeEmpty())
	})
})

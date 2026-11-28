package helpers

import (
	"fmt"
	"go/ast"
	"strings"
)

func annotated(doc *ast.CommentGroup) (bool, error) {
	if doc == nil { return false, nil }
	found := false
	for _, comment := range doc.List {
		line := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
		if line == "@PatchRequestDto" {
			if found { return false, fmt.Errorf("duplicate @PatchRequestDto") }
			found = true
		} else if strings.HasPrefix(line, "@PatchRequestDto ") {
			return false, fmt.Errorf("@PatchRequestDto does not take arguments")
		}
	}
	return found, nil
}

func isPatch(expr ast.Expr, aliases map[string]string) bool {
	index, ok := expr.(*ast.IndexExpr)
	if !ok { return false }
	selector, ok := index.X.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != "Patch" { return false }
	identifier, ok := selector.X.(*ast.Ident)
	return ok && aliases[identifier.Name] == patchImport
}

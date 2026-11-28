package helpers

import (
	"go/ast"
	"strings"
	"testing"
)

func TestDestinationTypeRejectsNonIdentifierQualifier(t *testing.T) {
	expr := &ast.SelectorExpr{X: &ast.SelectorExpr{X: ast.NewIdent("outer"), Sel: ast.NewIdent("inner")}, Sel: ast.NewIdent("Output")}
	_, err := destinationType("", "", nil, expr)
	if err == nil || !strings.Contains(err.Error(), "unsupported destination type") {
		t.Fatalf("expected unsupported destination type, got %v", err)
	}
}

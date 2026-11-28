package helpers

import (
	"go/ast"
)

func selectorType(expr ast.Expr, pkg,name string) bool {
	selector,ok:=expr.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name!=name{return false}
	base,ok:=selector.X.(*ast.Ident)
	return ok && base.Name==pkg
}

func identType(expr ast.Expr,name string)bool{
	id,ok:=expr.(*ast.Ident)
	return ok && id.Name==name
}

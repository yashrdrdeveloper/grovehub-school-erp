package helpers

import (
	"fmt"
	"go/ast"
)

// The explicit Search method is intentionally small: config and query
// mechanics remain generated, and callers see an ordinary method in impl.
func validateSearchDelegate(fn *ast.FuncDecl) error {
	if fn.Body == nil || len(fn.Body.List) != 1 {
		return fmt.Errorf("@Searchable Search must delegate to repository.searchGenerated(ctx, query)")
	}
	stmt, ok := fn.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(stmt.Results) != 1 {
		return fmt.Errorf("@Searchable Search must return repository.searchGenerated(ctx, query)")
	}
	call, ok := stmt.Results[0].(*ast.CallExpr)
	if !ok || len(call.Args) != 2 {
		return fmt.Errorf("@Searchable Search must call the generated search helper")
	}
	target, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || target.Sel.Name != "searchGenerated" { return fmt.Errorf("@Searchable Search must use searchGenerated") }
	owner, ok := target.X.(*ast.Ident)
	if !ok || len(fn.Recv.List[0].Names) != 1 || owner.Name != fn.Recv.List[0].Names[0].Name {
		return fmt.Errorf("@Searchable Search must delegate using its repository receiver")
	}
	if exprString(call.Args[0]) != "ctx" || exprString(call.Args[1]) != "query" {
		return fmt.Errorf("@Searchable Search must forward ctx and query unchanged")
	}
	return nil
}

func validateSignature(expr ast.Expr,entity string)error {
	fn,ok:=expr.(*ast.FuncType)
	if !ok||fn.Params==nil||fn.Results==nil||len(fn.Params.List)!=2||len(fn.Results.List)!=2 {
		return fmt.Errorf("requires (context.Context, coresearch.Query[string]) (coresearch.QueryResult[entity.%s], error)",entity)
	}
	if exprString(fn.Params.List[0].Type)!="context.Context" ||
		exprString(fn.Params.List[1].Type)!="coresearch.Query[string]" ||
		exprString(fn.Results.List[0].Type)!="coresearch.QueryResult[entity."+entity+"]" ||
		exprString(fn.Results.List[1].Type)!="error" {
		return fmt.Errorf("unsupported signature: expected context.Context, coresearch.Query[string] and QueryResult[entity.%s]",entity)
	}
	return nil
}

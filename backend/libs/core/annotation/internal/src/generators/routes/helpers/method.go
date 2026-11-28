package helpers

import (
	"fmt"
	"go/ast"
	"strconv"
	"strings"
)

func methodFor(doc *ast.CommentGroup) (route, error) {
	var result route
	if doc == nil {
		return result, nil
	}
	for _, comment := range doc.List {
		line := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
		if !strings.HasPrefix(line, "@") {
			continue
		}
		parts := strings.Fields(strings.TrimPrefix(line, "@"))
		if len(parts) == 0 || !verbs[parts[0]] {
			continue
		}
		if result.verb != "" {
			return result, fmt.Errorf("multiple HTTP verb annotations")
		}
		if len(parts) != 2 {
			return result, fmt.Errorf("@%s requires exactly one path", parts[0])
		}
		result.verb, result.path = parts[0], parts[1]
	}
	return result, nil
}

// routeLimitFor validates a limit of N requests per minute. Missing annotation
// means no route-specific limiter; it must never create a zero/negative limit.
func routeLimitFor(doc *ast.CommentGroup) (int, error) {
	if doc == nil {
		return 0, nil
	}
	found := false
	limit := 0
	for _, comment := range doc.List {
		line := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
		if line != "@RouteLimit" && !strings.HasPrefix(line, "@RouteLimit ") {
			continue
		}
		if found {
			return 0, fmt.Errorf("duplicate @RouteLimit annotation")
		}
		found = true
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return 0, fmt.Errorf("@RouteLimit requires positive requests per minute")
		}
		count, err := strconv.Atoi(fields[1])
		if err != nil || count < 1 || count > 100000 {
			return 0, fmt.Errorf("@RouteLimit must be 1..100000 requests per minute")
		}
		limit = count
	}
	return limit, nil
}

func validSignature(expr ast.Expr) bool {
	fn, ok := expr.(*ast.FuncType)
	if !ok || fn.Results != nil || fn.Params == nil || len(fn.Params.List) != 2 {
		return false
	}
	writer, ok := fn.Params.List[0].Type.(*ast.SelectorExpr)
	if !ok || writer.Sel.Name != "ResponseWriter" {
		return false
	}
	pkg, ok := writer.X.(*ast.Ident)
	if !ok || pkg.Name != "http" {
		return false
	}
	pointer, ok := fn.Params.List[1].Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	request, ok := pointer.X.(*ast.SelectorExpr)
	if !ok || request.Sel.Name != "Request" {
		return false
	}
	pkg, ok = request.X.(*ast.Ident)
	return ok && pkg.Name == "http"
}

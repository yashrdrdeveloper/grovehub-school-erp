package helpers

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
	"strconv"
	"strings"
)

func controllerPatterns(path string) (map[string]bool,error) {
	file,err := parser.ParseFile(token.NewFileSet(),path,nil,0)
	if err != nil { return nil,fmt.Errorf("read generated controller routes %s: %w",path,err) }
	patterns := map[string]bool{}
	ast.Inspect(file,func(node ast.Node)bool {
		call,ok := node.(*ast.CallExpr)
		if !ok || len(call.Args)<1 { return true }
		selector,ok := call.Fun.(*ast.SelectorExpr)
		if !ok || selector.Sel.Name!="Handle" { return true }
		value,ok := call.Args[0].(*ast.BasicLit)
		if !ok || value.Kind!=token.STRING { return true }
		pattern,err := strconv.Unquote(value.Value)
		if err==nil { patterns[pattern]=true }
		return true
	})
	return patterns,nil
}

// resolveRoute matches an exact registered endpoint or a service-prefix-free
// @RequestLine such as /protected/persons/{gid}. The generator derives the
// service prefix from the authoritative generated route registrations.
func resolveRoute(verb, declared string, patterns map[string]bool) (string,error) {
	if patterns[verb+" "+declared] { return declared,nil }
	var matches []string
	for pattern := range patterns {
		if !strings.HasPrefix(pattern,verb+" ") { continue }
		path := strings.TrimPrefix(pattern,verb+" ")
		if path != declared && strings.HasSuffix(path,declared) {
			matches = append(matches,path)
		}
	}
	if len(matches)==0 {
		return "",fmt.Errorf("@RequestLine %s %s has no matching generated controller route",verb,declared)
	}
	if len(matches)>1 {
		sort.Strings(matches)
		return "",fmt.Errorf("@RequestLine %s %s matches multiple routes; use an unambiguous path",verb,declared)
	}
	return matches[0],nil
}

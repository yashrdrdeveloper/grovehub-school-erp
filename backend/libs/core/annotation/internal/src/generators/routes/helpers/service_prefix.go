package helpers

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

// servicePrefix looks only at the service's RegisterHTTP function. Older
// service controllers without a service prefix retain their existing URLs.
func ServicePrefix(filename string) (string, bool, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return "", false, nil
	} else if err != nil {
		return "", false, err
	}
	file, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
	if err != nil {
		return "", false, err
	}
	var prefix string
	found := false
	for _, declaration := range file.Decls {
		fn, ok := declaration.(*ast.FuncDecl)
		if !ok || fn.Doc == nil {
			continue
		}
		value, annotated, err := prefixFor(fn.Doc)
		if err != nil {
			return "", false, fmt.Errorf("%s: %w", filename, err)
		}
		if !annotated {
			continue
		}
		if fn.Name.Name != "RegisterHTTP" || fn.Recv != nil {
			return "", false, fmt.Errorf("%s: service @Routes must annotate the RegisterHTTP function", filename)
		}
		if found {
			return "", false, fmt.Errorf("%s: duplicate service @Routes", filename)
		}
		if strings.ContainsAny(value, "{}*") {
			return "", false, fmt.Errorf("%s: service prefix cannot contain path variables or wildcards", filename)
		}
		found = true
		prefix = value
	}
	return prefix, found, nil
}

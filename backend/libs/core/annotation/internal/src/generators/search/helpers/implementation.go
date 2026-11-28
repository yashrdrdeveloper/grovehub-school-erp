package helpers

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"grovehub/libs/core/annotation/internal/src/helpers/layout"
	"strings"
)

func ValidateImplementation(dir string, specs []spec) error {
	files, err := layout.GoFiles(dir)
	if err != nil {
		return err
	}
	structs := make(map[string]bool)
	for _, filename := range files {
		file, err := parser.ParseFile(token.NewFileSet(), filename, nil, 0)
		if err != nil {
			return err
		}
		for _, decl := range file.Decls {
			switch node := decl.(type) {
			case *ast.GenDecl:
				if node.Tok != token.TYPE {
					continue
				}
				for _, typeSpec := range layout.TypeSpecs(node) {
					st, ok := typeSpec.Type.(*ast.StructType)
					if !ok {
						continue
					}
					for _, field := range st.Fields.List {
						for _, name := range field.Names {
							if name.Name == "db" {
								structs[typeSpec.Name.Name] = true
							}
						}
					}
				}
			case *ast.FuncDecl:
				if node.Recv == nil || node.Name.Name != "searchGenerated" {
					continue
				}
				for _, wanted := range specs {
					if strings.TrimPrefix(exprString(node.Recv.List[0].Type), "*") == wanted.impl {
						return fmt.Errorf("%s already implements searchGenerated; @Searchable owns that helper", wanted.impl)
					}
				}
			}
		}
	}
	for _, wanted := range specs {
		if !structs[wanted.impl] {
			return fmt.Errorf("%s needs a db field on its implementation struct", wanted.impl)
		}
	}
	return nil
}

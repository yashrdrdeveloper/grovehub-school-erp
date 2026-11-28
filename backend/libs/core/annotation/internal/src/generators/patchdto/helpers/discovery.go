package helpers

import (
	"go/token"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/helpers/layout"
	"fmt"
	"go/ast"
	"go/parser"
	"reflect"
	"strings"
)

func Discover(dir string) ([]patchType, string, error) {
	names, err := layout.GoFiles(dir)
	if err != nil { return nil, "", err }
	var specs []patchType
	seen := map[string]bool{}
	methods := map[string]bool{}
	var pkg string
	for _, filename := range names {
		if strings.HasSuffix(filename, "_test.go") || strings.HasSuffix(filename, ".gen.go") { continue }
		source,err:=layout.ParseSource(filename,parser.ParseComments)
		if err!=nil{return nil,"",err}
		file:=source.Tree
		if pkg == "" { pkg = file.Name.Name } else if file.Name.Name != pkg {
			return nil, "", fmt.Errorf("mixed packages in %s", dir)
		}
		aliases := make(map[string]string)
		for _,imp:=range source.Imports{aliases[imp.Name]=imp.Path}
		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Recv != nil && fn.Name.Name == "MarshalJSON" {
				recv := fn.Recv.List[0].Type
				if ptr, ok := recv.(*ast.StarExpr); ok { recv = ptr.X }
				if id, ok := recv.(*ast.Ident); ok { methods[id.Name] = true }
			}
			group, ok := decl.(*ast.GenDecl)
			if !ok || group.Tok != token.TYPE { continue }
			for _,spec:=range layout.TypeSpecs(group){
				doc := spec.Doc
				if doc == nil { doc = group.Doc }
				found, err := annotated(doc)
				if err != nil { return nil, "", fmt.Errorf("%s.%s: %w", filename, spec.Name.Name, err) }
				if !found { continue }
				if seen[spec.Name.Name] { return nil, "", fmt.Errorf("duplicate @PatchRequestDto: %s", spec.Name.Name) }
				seen[spec.Name.Name] = true
				st, ok := spec.Type.(*ast.StructType)
				if !ok || spec.TypeParams != nil { return nil, "", fmt.Errorf("%s: @PatchRequestDto requires a non-generic struct", spec.Name.Name) }
				result := patchType{name: spec.Name.Name}
				used := make(map[string]bool)
				for _, field := range st.Fields.List {
					if len(field.Names) != 1 || !field.Names[0].IsExported() {
						return nil, "", fmt.Errorf("%s: fields must be exported, named and ungrouped", result.name)
					}
					if !isPatch(field.Type, aliases) {
						return nil, "", fmt.Errorf("%s.%s: must use %s.Patch[T]", result.name, field.Names[0].Name, patchImport)
					}
					if field.Tag == nil { return nil, "", fmt.Errorf("%s.%s: missing json tag", result.name, field.Names[0].Name) }
					raw:=source.Tags[field]
					jsonName := strings.Split(reflect.StructTag(raw).Get("json"), ",")[0]
					if jsonName == "" || jsonName == "-" || used[jsonName] {
						return nil, "", fmt.Errorf("%s.%s: invalid or duplicate json field %q", result.name, field.Names[0].Name, jsonName)
					}
					used[jsonName] = true
					result.fields = append(result.fields, patchField{name: field.Names[0].Name, jsonName: jsonName})
				}
				specs = append(specs, result)
			}
		}
	}
	for _, spec := range specs {
		if methods[spec.name] { return nil, "", fmt.Errorf("%s: remove handwritten MarshalJSON; @PatchRequestDto generates it", spec.name) }
	}
	return specs, pkg, nil
}

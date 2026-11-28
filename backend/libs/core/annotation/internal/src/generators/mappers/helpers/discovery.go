package helpers

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"grovehub/libs/core/annotation/internal/src/helpers/layout"
	"strings"
)

func Discover(backend, modulePath, dir string) (mappingGroup, bool, error) {
	return DiscoverWithRenderer(backend, modulePath, dir, layout.Render)
}

func DiscoverWithRenderer(backend, modulePath, dir string, render layout.Renderer) (mappingGroup, bool, error) {
	files, err := layout.GoFiles(dir)
	if err != nil {
		return mappingGroup{}, false, err
	}
	var result mappingGroup
	found := false
	for _, filePath := range files {
		if strings.HasSuffix(filePath, "_test.go") {
			continue
		}
		source, err := layout.ParseSource(filePath, parser.ParseComments)
		if err != nil {
			return result, false, err
		}
		file, set := source.Tree, source.Set
		for _, decl := range file.Decls {
			group, ok := decl.(*ast.GenDecl)
			if !ok || group.Tok != token.TYPE {
				continue
			}
			for _, spec := range layout.TypeSpecs(group) {
				iface, ok := spec.Type.(*ast.InterfaceType)
				if !ok {
					continue
				}
				doc := spec.Doc
				if doc == nil {
					doc = group.Doc
				}
				receiver, present, err := argument(doc, "Mapper")
				if err != nil {
					return result, false, err
				}
				if !present {
					continue
				}
				if found {
					return result, false, fmt.Errorf("one @Mapper interface per service is supported")
				}
				if !token.IsIdentifier(receiver) {
					return result, false, fmt.Errorf("invalid @Mapper receiver %q", receiver)
				}
				found = true
				result.receiver = receiver
				result.imports = source.Imports
				for _, member := range iface.Methods.List {
					if len(member.Names) != 1 {
						return result, false, fmt.Errorf("embedded mapper interfaces are not supported")
					}
					rules, err := mappingRules(member.Doc)
					if err != nil {
						return result, false, fmt.Errorf("%s: %w", member.Names[0].Name, err)
					}
					if len(rules) == 0 {
						continue
					}
					fn, ok := member.Type.(*ast.FuncType)
					if !ok || fn.Params == nil || len(fn.Params.List) != 1 || len(fn.Params.List[0].Names) != 1 || fn.Results == nil || len(fn.Results.List) != 1 {
						return result, false, fmt.Errorf("%s: @Mapping requires one named input and one struct result", member.Names[0].Name)
					}
					destination, err := destinationType(backend, modulePath, source.Imports, fn.Results.List[0].Type)
					if err != nil {
						return result, false, fmt.Errorf("%s: %w", member.Names[0].Name, err)
					}
					if len(rules) != len(destination) {
						return result, false, fmt.Errorf("%s: mapped %d fields, destination requires %d", member.Names[0].Name, len(rules), len(destination))
					}
					ordered := make([]fieldMapping, 0, len(destination))
					for _, field := range destination {
						expr, ok := rules[field]
						if !ok {
							return result, false, fmt.Errorf("%s: unmapped destination field %s", member.Names[0].Name, field)
						}
						ordered = append(ordered, fieldMapping{name: field, expression: expr})
					}
					signature, err := render(set, fn)
					if err != nil {
						return result, false, err
					}
					resultType, err := render(set, fn.Results.List[0].Type)
					if err != nil {
						return result, false, err
					}
					result.methods = append(result.methods, mappingMethod{
						name: member.Names[0].Name, signature: strings.TrimPrefix(signature, "func"),
						output: resultType, fields: ordered,
					})
				}
			}
		}
	}
	return result, found, nil
}

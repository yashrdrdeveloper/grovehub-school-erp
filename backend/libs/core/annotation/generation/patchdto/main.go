// Command patchdto generates JSON marshal methods for @PatchRequestDto
// structs into the existing service-local, ignored Go source overlay.
package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

const patchImport = "grovehub/libs/core/annotation/dto"

type patchField struct {
	name string
	jsonName string
}

type patchType struct {
	name string
	fields []patchField
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: patchdto <backend-directory>")
		os.Exit(2)
	}
	if err := generateAll(os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Every DTO package is discovered independently under dto/, allowing
// resource-owned DTOs (dto/person) and legacy request/ packages alike.
func generateAll(backend string) error {
	services, err := os.ReadDir(filepath.Join(backend, "services"))
	if err != nil { return err }
	for _, entry := range services {
		if !entry.IsDir() { continue }
		owner := filepath.Join(backend, "services", entry.Name())
		roots := []string{filepath.Join(owner, "dto"),filepath.Join(owner, "client")}
        for _,root := range roots {
            if _, err := os.Stat(root); os.IsNotExist(err) { continue } else if err != nil { return err }
            err = filepath.WalkDir(root, func(dir string, item os.DirEntry, walkErr error) error {
			if walkErr != nil { return walkErr }
			if !item.IsDir() { return nil }
			if item.Name() == ".generated" || item.Name() == "vendor" { return filepath.SkipDir }
			specs, pkg, err := discover(dir)
			if err != nil { return fmt.Errorf("%s: %w", dir, err) }
			rel, err := filepath.Rel(owner, dir)
			if err != nil { return err }
			output := filepath.Join(owner, ".generated", rel, "zz_patchdto.gen.go")
			if len(specs) == 0 {
				if err := os.Remove(output); err != nil && !os.IsNotExist(err) { return err }
				return nil
			}
			return write(output, pkg, specs)
		})
		if err != nil { return err }
        }
	}
	return nil
}

func discover(dir string) ([]patchType, string, error) {
	names, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil { return nil, "", err }
	sort.Strings(names)
	var specs []patchType
	seen := map[string]bool{}
	methods := map[string]bool{}
	var pkg string
	for _, filename := range names {
		if strings.HasSuffix(filename, "_test.go") || strings.HasSuffix(filename, ".gen.go") { continue }
		file, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
		if err != nil { return nil, "", err }
		if pkg == "" { pkg = file.Name.Name } else if file.Name.Name != pkg {
			return nil, "", fmt.Errorf("mixed packages in %s", dir)
		}
		aliases := make(map[string]string)
		for _, imp := range file.Imports {
			name, err := strconv.Unquote(imp.Path.Value)
			if err != nil { return nil, "", err }
			alias := path.Base(name)
			if imp.Name != nil { alias = imp.Name.Name }
			aliases[alias] = name
		}
		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok && fn.Recv != nil && fn.Name.Name == "MarshalJSON" {
				recv := fn.Recv.List[0].Type
				if ptr, ok := recv.(*ast.StarExpr); ok { recv = ptr.X }
				if id, ok := recv.(*ast.Ident); ok { methods[id.Name] = true }
			}
			group, ok := decl.(*ast.GenDecl)
			if !ok || group.Tok != token.TYPE { continue }
			for _, node := range group.Specs {
				spec, ok := node.(*ast.TypeSpec)
				if !ok { continue }
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
					raw, err := strconv.Unquote(field.Tag.Value)
					if err != nil { return nil, "", err }
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

func write(output, pkg string, specs []patchType) error {
	var b strings.Builder
	fmt.Fprintf(&b, "// Code generated by GrooveHub @PatchRequestDto. DO NOT EDIT.\npackage %s\n\nimport \"encoding/json\"\n\n", pkg)
	for _, spec := range specs {
		fmt.Fprintf(&b, "func (request %s) MarshalJSON() ([]byte, error) {\n", spec.name)
		fmt.Fprintf(&b, "fields := make(map[string]any, %d)\n", len(spec.fields))
		for _, field := range spec.fields {
			fmt.Fprintf(&b, "if request.%s.Set { fields[%q] = request.%s.Value }\n", field.name, field.jsonName, field.name)
		}
		b.WriteString("return json.Marshal(fields)\n}\n\n")
	}
	source, err := format.Source([]byte(b.String()))
	if err != nil { return err }
	if old, err := os.ReadFile(output); err == nil && bytes.Equal(old, source) { return nil } else if err != nil && !os.IsNotExist(err) { return err }
	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil { return err }
	return os.WriteFile(output, source, 0644)
}

package helpers

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"sort"
	"strings"
)

// @Searchable belongs to the concrete Search method, not the public interface.
// The developer-maintained Search method delegates to generated searchGenerated.
func Discover(dir, modulePath, service string) ([]spec, error) {
	interfaces, err := repositoryInterfaces(dir)
	if err != nil { return nil, err }
	paths, err := filepath.Glob(filepath.Join(dir, "impl", "*.go"))
	if err != nil { return nil, err }
	sort.Strings(paths)
	var result []spec
	seen := make(map[string]bool)
	for _, filename := range paths {
		if strings.HasSuffix(filename, "_test.go") { continue }
		file, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
		if err != nil { return nil, err }
		for _, decl := range file.Decls {
			method, ok := decl.(*ast.FuncDecl)
			if !ok || !hasDirective(method.Doc, "Searchable") { continue }
			if method.Recv == nil || len(method.Recv.List) != 1 || method.Name.Name != "Search" {
				return nil, fmt.Errorf("%s: @Searchable requires a concrete Search method", filename)
			}
			receiverType := strings.TrimPrefix(exprString(method.Recv.List[0].Type), "*")
			if !token.IsIdentifier(receiverType) {
				return nil, fmt.Errorf("%s: @Searchable requires a named repository receiver", filename)
			}
			if seen[receiverType] {
				return nil, fmt.Errorf("%s: duplicate @Searchable Search for %s", filename, receiverType)
			}
			seen[receiverType] = true
			entityName, err := repositoryEntity(interfaces, receiverType)
			if err != nil { return nil, fmt.Errorf("%s.Search: %w", receiverType, err) }
			if err := validateSignature(method.Type, entityName); err != nil {
				return nil, fmt.Errorf("%s.Search: %w", receiverType, err)
			}
			if err := validateSearchDelegate(method); err != nil {
				return nil, fmt.Errorf("%s.Search: %w", receiverType, err)
			}
			opts, err := searchOptions(directive(method.Doc, "Searchable"))
			if err != nil { return nil, fmt.Errorf("%s.Search: %w", receiverType, err) }
			if opts["entity"] != "" && opts["entity"] != entityName {
				return nil, fmt.Errorf("%s.Search: @Searchable entity must match @Repository", receiverType)
			}
			if opts["impl"] != "" && opts["impl"] != receiverType {
				return nil, fmt.Errorf("%s.Search: @Searchable impl must match concrete receiver", receiverType)
			}
			filterFields, err := optionalFields(method.Doc, "Filterable")
			if err != nil { return nil, err }
			sortFields, err := optionalFields(method.Doc, "Sortable")
			if err != nil { return nil, err }
			cursorName := ""
			if raw := directive(method.Doc, "Cursor"); raw != "" {
				cursorArgs, err := options(raw, []string{"field"})
				if err != nil { return nil, fmt.Errorf("@Cursor: %w", err) }
				cursorName = cursorArgs["field"]
				if len(sortFields) > 0 && !contains(sortFields, cursorName) {
					return nil, fmt.Errorf("@Cursor field must be sortable")
				}
			}
			configDir, configImport, configAlias, configName, err := resolveConfig(dir, modulePath, service, opts["config"])
			if err != nil { return nil, err }
			if err := validateConfig(configDir, configName, filterFields, sortFields, cursorName); err != nil { return nil, err }
			result = append(result, spec{
				service: service, impl: receiverType, entity: entityName, config: opts["config"],
				alias: opts["alias"], configImport: configImport, configAlias: configAlias,
			})
		}
	}
	return result, nil
}

// Reject misplaced @Searchable annotations before emitting a search helper.
func repositoryInterfaces(dir string) ([]*ast.File,error) {
	paths, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil { return nil,err }
	var files []*ast.File
	for _, filename := range paths {
		if strings.HasSuffix(filename, "_test.go") { continue }
		file, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
		if err != nil { return nil,err }
		files=append(files,file)
		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE { continue }
			for _, node := range gen.Specs {
				t, ok := node.(*ast.TypeSpec)
				if !ok { continue }
				iface, ok := t.Type.(*ast.InterfaceType)
				if !ok { continue }
				for _, member := range iface.Methods.List {
					if hasDirective(member.Doc, "Searchable") {
						return nil,fmt.Errorf("%s: move @Searchable from %s interface to the concrete Search method", filename, t.Name.Name)
					}
				}
			}
		}
	}
	return files,nil
}

// repositoryEntity ensures the annotated implementation satisfies the
// explicit repository interface and its @Repository ownership declaration.
func repositoryEntity(files []*ast.File, receiver string) (string, error) {
	for _, file := range files {
		for _, decl := range file.Decls {
			group, ok := decl.(*ast.GenDecl)
			if !ok || group.Tok != token.TYPE { continue }
			for _, node := range group.Specs {
				typ, ok := node.(*ast.TypeSpec)
				if !ok || typ.Name.Name != receiver { continue }
				iface, ok := typ.Type.(*ast.InterfaceType)
				if !ok { continue }
				doc := typ.Doc
				if doc == nil { doc = group.Doc }
				repoOpts, err := options(directive(doc, "Repository"), []string{"entity"})
				if err != nil { return "", fmt.Errorf("matching @Repository interface required: %w", err) }
				entity := repoOpts["entity"]
				for _, member := range iface.Methods.List {
					if len(member.Names) != 1 { continue }
					if member.Names[0].Name != "Search" { continue }
					if hasDirective(member.Doc, "Searchable") {
						return "", fmt.Errorf("move @Searchable from interface to %s.Search implementation", receiver)
					}
					if err := validateSignature(member.Type, entity); err != nil { return "", err }
					return entity, nil
				}
				return "", fmt.Errorf("repository interface %s must declare Search", receiver)
			}
		}
	}
	return "", fmt.Errorf("missing @Repository interface for concrete %s", receiver)
}

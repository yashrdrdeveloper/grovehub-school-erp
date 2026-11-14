// Command routes generates typed net/http registrations from annotated
// controller interfaces. The service-level @Routes declaration lives on
// routes.RegisterHTTP; controller package paths declare access tiers.
//
// This generator only groups routes. No tier enforces authentication,
// authorization or internal-network restrictions.
package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type route struct {
	method string
	verb   string
	path   string
}

type routeGroup struct {
	name   string
	tier   string
	resource string
	routes []route
}

var verbs = map[string]bool{
	"GET": true, "POST": true, "PUT": true,
	"PATCH": true, "DELETE": true,
}

var accessTiers = []string{"public", "protected", "restricted"}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: routes <backend-dir>")
		os.Exit(2)
	}
	if err := generateAll(os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func generateAll(backend string) error {
	services, err := os.ReadDir(filepath.Join(backend, "services"))
	if err != nil {
		return err
	}
	for _, service := range services {
		if !service.IsDir() {
			continue
		}
		name := service.Name()
		owner := filepath.Join(backend, "services", name)
		controllerDir := filepath.Join(owner, "controller")
        if info, err := os.Stat(filepath.Join(owner, "internal", "app", "controller")); err == nil && info.IsDir() {
            controllerDir = filepath.Join(owner, "internal", "app", "controller")
        }
		if info, err := os.Stat(controllerDir); os.IsNotExist(err) {
			continue
		} else if err != nil {
			return err
		} else if !info.IsDir() {
			return fmt.Errorf("%s: controller must be a directory", name)
		}

		routesDir := filepath.Join(filepath.Dir(controllerDir), "routes")
        servicePath, configured, err := servicePrefix(filepath.Join(routesDir, "http.go"))
		if err != nil {
			return fmt.Errorf("%s routes: %w", name, err)
		}
		groups, err := discover(controllerDir, servicePath, configured)
		if err != nil {
			return fmt.Errorf("%s routes: %w", name, err)
		}
		rel, err := filepath.Rel(owner, routesDir)
        if err != nil { return err }
        output := filepath.Join(owner, ".generated", rel, "zz_routes.gen.go")
		if len(groups) == 0 {
			if err := os.Remove(output); err != nil && !os.IsNotExist(err) {
				return err
			}
			continue
		}
		if err := writeRoutes(name, groups, output); err != nil {
			return fmt.Errorf("%s routes: %w", name, err)
		}
	}
	return nil
}

// servicePrefix looks only at the service's RegisterHTTP function. Older
// service controllers without a service prefix retain their existing URLs.
func servicePrefix(filename string) (string, bool, error) {
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

func discover(dir, servicePath string, configured bool) ([]routeGroup, error) {
	// Support both legacy controller/<tier> and resource-owned
	// controller/<resource>/<tier> packages.
	type category struct{ dir, tier, resource string }
	directories := []category{{dir:dir}}
	for _, tier := range accessTiers {
		directories = append(directories, category{dir:filepath.Join(dir,tier),tier:tier})
	}
	resources, err := os.ReadDir(dir)
	if err != nil { return nil, err }
	for _, resource := range resources {
		if !resource.IsDir() { continue }
		if containsTier(resource.Name()) { continue }
		for _, tier := range accessTiers {
			directories=append(directories,category{dir:filepath.Join(dir,resource.Name(),tier),tier:tier,resource:resource.Name()})
		}
	}

	groups := make([]routeGroup, 0)
	patterns := make(map[string]string)
	for _, directory := range directories {
		info, err := os.Stat(directory.dir)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			return nil, fmt.Errorf("%s: controller category must be a directory", directory.dir)
		}
		files, err := filepath.Glob(filepath.Join(directory.dir, "*.go"))
		if err != nil {
			return nil, err
		}
		sort.Strings(files)
		for _, filename := range files {
			if strings.HasSuffix(filename, "_test.go") {
				continue
			}
			file, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
			if err != nil {
				return nil, err
			}
			wantPackage := directory.tier
			if wantPackage == "" {
				wantPackage = "controller"
			}
			if file.Name.Name != wantPackage {
				return nil, fmt.Errorf("%s: expected controller package %q", filename, wantPackage)
			}
			for _, declaration := range file.Decls {
				gen, ok := declaration.(*ast.GenDecl)
				if !ok || gen.Tok != token.TYPE {
					continue
				}
				for _, node := range gen.Specs {
					typ, ok := node.(*ast.TypeSpec)
					if !ok {
						continue
					}
					doc := typ.Doc
					if doc == nil {
						doc = gen.Doc
					}
					resourcePath, annotated, err := prefixFor(doc)
					if err != nil {
						return nil, fmt.Errorf("%s.%s: %w", filename, typ.Name.Name, err)
					}
					if !annotated {
						continue
					}
					if configured && directory.tier == "" {
						return nil, fmt.Errorf("%s: controller @Routes must live in public/, protected/ or restricted/ when service @Routes is configured", filename)
					}
					if !configured && directory.tier != "" {
						return nil, fmt.Errorf("%s: categorized controller requires a service @Routes prefix on routes.RegisterHTTP", filename)
					}
					iface, ok := typ.Type.(*ast.InterfaceType)
					if !ok || !strings.HasSuffix(typ.Name.Name, "Controller") {
						return nil, fmt.Errorf("%s: @Routes requires a Controller interface", filename)
					}
					group := routeGroup{name: typ.Name.Name, tier: directory.tier, resource: directory.resource}
					for _, member := range iface.Methods.List {
						if len(member.Names) != 1 {
							return nil, fmt.Errorf("%s cannot embed interfaces", group.name)
						}
						if !validSignature(member.Type) {
							return nil, fmt.Errorf("%s.%s must accept http.ResponseWriter and *http.Request", group.name, member.Names[0].Name)
						}
						item, err := methodFor(member.Doc)
						if err != nil {
							return nil, err
						}
						if item.verb == "" {
							return nil, fmt.Errorf("%s.%s missing @GET/@POST/etc", group.name, member.Names[0].Name)
						}
						item.method = member.Names[0].Name
						if err := validatePath(item.path); err != nil {
							return nil, fmt.Errorf("%s.%s: %w", group.name, item.method, err)
						}
						item.path = joinPaths(servicePath, directory.tier, resourcePath, item.path)
						key := item.verb + " " + item.path
						if owner, duplicate := patterns[key]; duplicate {
							return nil, fmt.Errorf("duplicate route %s: %s and %s.%s", key, owner, group.name, item.method)
						}
						patterns[key] = group.name + "." + item.method
						group.routes = append(group.routes, item)
					}
					groups = append(groups, group)
				}
			}
		}
	}
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].tier+":"+groups[i].name < groups[j].tier+":"+groups[j].name
	})
	return groups, nil
}

func containsTier(value string) bool {
	for _, tier := range accessTiers { if tier==value { return true } }
	return false
}

func prefixFor(doc *ast.CommentGroup) (string, bool, error) {
	if doc == nil {
		return "", false, nil
	}
	found := false
	var prefix string
	for _, comment := range doc.List {
		line := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
		if line != "@Routes" && !strings.HasPrefix(line, "@Routes ") {
			continue
		}
		if found {
			return "", true, fmt.Errorf("duplicate @Routes declaration")
		}
		found = true
		prefix = strings.TrimSpace(strings.TrimPrefix(line, "@Routes"))
		if err := validatePath(prefix); err != nil {
			return "", true, err
		}
	}
	return prefix, found, nil
}

func validatePath(value string) error {
	if !strings.HasPrefix(value, "/") || strings.Contains(value, "//") ||
		strings.ContainsAny(value, " \t\r\n?#") {
		return fmt.Errorf("invalid route path %q", value)
	}
	for _, part := range strings.Split(value, "/") {
		if part == "." || part == ".." {
			return fmt.Errorf("invalid route path %q", value)
		}
	}
	return nil
}

func joinPaths(parts ...string) string {
	segments := make([]string, 0, len(parts))
	for _, part := range parts {
		if segment := strings.Trim(part, "/"); segment != "" {
			segments = append(segments, segment)
		}
	}
	return "/" + strings.Join(segments, "/")
}

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

func generatedRegistrationName(group routeGroup) string {
	tier := ""
	switch group.tier {
	case "public":
		tier = "Public"
	case "protected":
		tier = "Protected"
	case "restricted":
		tier = "Restricted"
	}
	return "register" + tier + strings.TrimSuffix(group.name, "Controller") + "Routes"
}

func groupImportAlias(group routeGroup) string {
	pkg:=group.tier
	if pkg=="" {pkg="controller"}
	if group.resource!="" {return group.resource+pkg}
	return pkg
}

func writeRoutes(service string, groups []routeGroup, output string) error {
	var b strings.Builder
	b.WriteString("// Code generated by GrooveHub @Routes. DO NOT EDIT.\n")
	b.WriteString("// Route tiers are URL categories only; security enforcement is a separate gateway responsibility.\n")
	b.WriteString("package routes\n\nimport (\n\"net/http\"\n")
	b.WriteString("httptracing \"grovehub/libs/core/http/tracing\"\n")
	imports := map[string]string{}
	for _, group := range groups {
		alias := groupImportAlias(group)
		importPath := "grovehub/services/" + service + "/controller"
        if strings.Contains(filepath.ToSlash(output), "/internal/app/routes/") {
            importPath = "grovehub/services/" + service + "/internal/app/controller"
        }
		if group.resource != "" { importPath += "/" + group.resource }
		if group.tier != "" { importPath += "/" + group.tier }
		imports[alias] = importPath
	}
	aliases := make([]string, 0, len(imports))
	for alias := range imports {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)
	for _, alias := range aliases {
		fmt.Fprintf(&b, "%s %q\n", alias, imports[alias])
	}
	b.WriteString(")\n\n")
	for _, group := range groups {
		fmt.Fprintf(&b, "func %s(mux *http.ServeMux, c %s.%s) {\n",
			generatedRegistrationName(group), groupImportAlias(group), group.name)
		for _, item := range group.routes {
			pattern := item.verb + " " + item.path
			fmt.Fprintf(&b, "mux.Handle(%q, httptracing.ServiceHandler(%q, %q, http.HandlerFunc(c.%s)))\n",
				pattern, service, pattern, item.method)
		}
		b.WriteString("}\n\n")
	}
	source, err := format.Source([]byte(b.String()))
	if err != nil {
		return fmt.Errorf("format generated routes: %w", err)
	}
	if old, err := os.ReadFile(output); err == nil && bytes.Equal(old, source) {
		return nil
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return err
	}
	return os.WriteFile(output, source, 0644)
}

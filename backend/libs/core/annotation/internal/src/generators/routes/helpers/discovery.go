package helpers

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func discover(dir, servicePath string, configured bool) ([]routeGroup, error) {
	return DiscoverWithTier(dir, servicePath, configured, false)
}

func DiscoverWithTier(dir, servicePath string, configured, restrictedOnly bool) ([]routeGroup, error) {
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
						limit, err := routeLimitFor(member.Doc)
						if err != nil {
							return nil, fmt.Errorf("%s.%s: %w", group.name, member.Names[0].Name, err)
						}
						item.limit = limit
						item.method = member.Names[0].Name
						if err := validatePath(item.path); err != nil {
							return nil, fmt.Errorf("%s.%s: %w", group.name, item.method, err)
						}
						routeTier := directory.tier
						if restrictedOnly {
							if !configured || routeTier == "" { return nil, fmt.Errorf("%s: restricted project requires categorized service routes", filename) }
							routeTier = "restricted"
						}
						item.path = joinPaths(servicePath, routeTier, resourcePath, item.path)
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

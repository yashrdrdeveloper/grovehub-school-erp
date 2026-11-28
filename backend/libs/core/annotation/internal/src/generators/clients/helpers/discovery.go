package helpers

import (
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/helpers/layout"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"strconv"
	"strings"
)

func Discover(dir, generatedRoutes string) ([]clientProxy,error) {
	return DiscoverWithRenderer(dir,generatedRoutes,layout.Render)
}

func DiscoverWithRenderer(dir,generatedRoutes string,render layout.Renderer)([]clientProxy,error){
	files,err := layout.GoFiles(dir)
	if err != nil { return nil,err }
	var result []clientProxy
	var patterns map[string]bool
	for _,path := range files {
		if strings.HasSuffix(path,"_test.go") { continue }
		set := token.NewFileSet()
		file,err := parser.ParseFile(set,path,nil,parser.ParseComments)
		if err != nil { return nil,err }
		for _,decl := range file.Decls {
			group,ok := decl.(*ast.GenDecl)
			if !ok || group.Tok != token.TYPE { continue }
			for _, spec := range layout.TypeSpecs(group) {
				iface,ok := spec.Type.(*ast.InterfaceType)
				if !ok { continue }
				doc := spec.Doc
				if doc == nil { doc = group.Doc }
				receiver,found,err := directive(doc,"Client")
				if err != nil { return nil,err }
				if !found { continue }
				if !token.IsIdentifier(receiver) { return nil,fmt.Errorf("%s: invalid @Client receiver %q",spec.Name.Name,receiver) }
				proxy := clientProxy{receiver:receiver,imports:signatureImports(file.Imports,iface)}
				for _,member := range iface.Methods.List {
					if len(member.Names)!=1 { return nil,fmt.Errorf("annotated @Client %s cannot embed interfaces",spec.Name.Name) }
					name := spec.Name.Name+"."+member.Names[0].Name
					target,hasProxy,err := directive(member.Doc,"Proxy")
					if err != nil { return nil,err }
					line,hasLine,err := requestLine(member.Doc)
					if err != nil { return nil,fmt.Errorf("%s: %w",name,err) }
					if hasProxy==hasLine { return nil,fmt.Errorf("%s: declare exactly one of @Proxy or @RequestLine",name) }
					// A successfully parsed named interface method always has a function type.
					fn := member.Type.(*ast.FuncType)
					var args []string
					for _,param := range fn.Params.List {
						if len(param.Names)==0 { return nil,fmt.Errorf("%s: all arguments must be named",name) }
						for _,arg := range param.Names { args=append(args,arg.Name) }
					}
					signature,err:=render(set,fn)
					if err!=nil{return nil,err}
					method := proxyMethod{name:member.Names[0].Name,signature:strings.TrimPrefix(signature,"func"),
						arguments:strings.Join(args,", "),hasResults:fn.Results!=nil&&len(fn.Results.List)>0}
					if hasProxy {
						parts := strings.Split(target,".")
						if len(parts)!=2 || !token.IsIdentifier(parts[0]) || !token.IsIdentifier(parts[1]) {
							return nil,fmt.Errorf("%s: invalid @Proxy binding %q",name,target)
						}
						method.delegate = target
					} else {
						if patterns==nil {
							patterns,err = controllerPatterns(generatedRoutes)
							if err != nil { return nil,err }
						}
						route,err := resolveRoute(line.verb, line.path, patterns)
						if err != nil { return nil,fmt.Errorf("%s: %w",name,err) }
						line.path = route
						ep,err := endpointArguments(set,fn,args,line,render)
						if err != nil { return nil,fmt.Errorf("%s: %w",name,err) }
						method.endpoint = &ep
					}
					proxy.methods = append(proxy.methods,method)
				}
				result=append(result,proxy)
			}
		}
	}
	return result,nil
}

// signatureImports retains only the imports needed by the generated method
// signatures. Facade helpers, constants, and initialization belong to the
// resource package and must not be copied into its implementation package.
func signatureImports(imports []*ast.ImportSpec, iface *ast.InterfaceType) []*ast.ImportSpec {
	used := map[string]bool{}
	ast.Inspect(iface, func(node ast.Node) bool {
		if selector, ok := node.(*ast.SelectorExpr); ok {
			if qualifier, ok := selector.X.(*ast.Ident); ok {
				used[qualifier.Name] = true
			}
		}
		return true
	})
	var result []*ast.ImportSpec
	for _, spec := range imports {
		name := path.Base(strings.Trim(spec.Path.Value, "\"\x60"))
		if spec.Name != nil {
			name = spec.Name.Name
		} else if strings.HasPrefix(name, "v") {
			// Semantic import versions use the preceding path component as the
			// conventional package name (example/module/v2 -> module).
			if _, err := strconv.Atoi(strings.TrimPrefix(name, "v")); err == nil {
				name = path.Base(path.Dir(strings.Trim(spec.Path.Value, "\"\x60")))
			}
		}
		// Dot imports cannot be inferred from selectors. Preserve them; blank
		// imports are facade initialization, never part of a method signature.
		if name == "." || used[name] {
			result = append(result, spec)
		}
	}
	return result
}

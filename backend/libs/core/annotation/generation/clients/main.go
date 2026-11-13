// Command clients generates typed client methods from @Client contracts.
//
// @RequestLine calls the verified generated controller route in-process,
// preserving the HTTP binding, validation and response pipeline.
// Legacy @Proxy methods remain supported for other service clients.
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
	"strconv"
	"strings"
)

type endpointMethod struct {
	verb, path, resultType string
	created bool
}

type proxyMethod struct {
	name, signature, delegate, arguments string
	hasResults bool
	endpoint *endpointMethod
}

type clientProxy struct {
	receiver string
	imports []*ast.ImportSpec
	methods []proxyMethod
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: clients <backend-dir>")
		os.Exit(2)
	}
	if err := generateAll(os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func generateAll(backend string) error {
	services, err := os.ReadDir(filepath.Join(backend, "services"))
	if err != nil { return err }
	for _,service := range services {
		if !service.IsDir() { continue }
		owner := filepath.Join(backend,"services",service.Name())
		clientDir := filepath.Join(owner,"client")
		if info,err := os.Stat(clientDir); err != nil || !info.IsDir() { continue }
		// Discover both legacy flat clients and client/<resource> contracts.
		entries,err:=os.ReadDir(clientDir)
		if err!=nil {return err}
		directories:=[]string{clientDir}
		for _,entry:=range entries {
			if entry.IsDir() && entry.Name()!="impl" && !strings.HasPrefix(entry.Name(),".") {
				directories=append(directories,filepath.Join(clientDir,entry.Name()))
			}
		}
		for _,dir:=range directories {
			proxies,err:=discover(dir,filepath.Join(owner,".generated","routes","zz_routes.gen.go"))
			if err!=nil{return fmt.Errorf("%s %s: %w",service.Name(),dir,err)}
			rel,err:=filepath.Rel(clientDir,dir)
			if err!=nil{return err}
			outDir:=filepath.Join(owner,".generated","client")
			if rel!="." {outDir=filepath.Join(outDir,rel)}
			output:=filepath.Join(outDir,"impl","zz_clients.gen.go")
			if len(proxies)==0 {
				if err:=os.Remove(output);err!=nil&&!os.IsNotExist(err){return err}
				continue
			}
			if err:=writeProxies(proxies,output);err!=nil{return fmt.Errorf("%s %s: %w",service.Name(),dir,err)}
		}
	}
	return nil
}

func discover(dir, generatedRoutes string) ([]clientProxy,error) {
	files,err := filepath.Glob(filepath.Join(dir,"*.go"))
	if err != nil { return nil,err }
	sort.Strings(files)
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
			for _,node := range group.Specs {
				spec,ok := node.(*ast.TypeSpec)
				if !ok { continue }
				iface,ok := spec.Type.(*ast.InterfaceType)
				if !ok { continue }
				doc := spec.Doc
				if doc == nil { doc = group.Doc }
				receiver,found,err := directive(doc,"Client")
				if err != nil { return nil,err }
				if !found { continue }
				if !token.IsIdentifier(receiver) { return nil,fmt.Errorf("%s: invalid @Client receiver %q",spec.Name.Name,receiver) }
				proxy := clientProxy{receiver:receiver,imports:file.Imports}
				for _,member := range iface.Methods.List {
					if len(member.Names)!=1 { return nil,fmt.Errorf("annotated @Client %s cannot embed interfaces",spec.Name.Name) }
					name := spec.Name.Name+"."+member.Names[0].Name
					target,hasProxy,err := directive(member.Doc,"Proxy")
					if err != nil { return nil,err }
					line,hasLine,err := requestLine(member.Doc)
					if err != nil { return nil,fmt.Errorf("%s: %w",name,err) }
					if hasProxy==hasLine { return nil,fmt.Errorf("%s: declare exactly one of @Proxy or @RequestLine",name) }
					fn,ok := member.Type.(*ast.FuncType)
					if !ok { return nil,fmt.Errorf("%s: expected function type",name) }
					var args []string
					for _,param := range fn.Params.List {
						if len(param.Names)==0 { return nil,fmt.Errorf("%s: all arguments must be named",name) }
						for _,arg := range param.Names { args=append(args,arg.Name) }
					}
					var signature bytes.Buffer
					if err := format.Node(&signature,set,fn); err != nil { return nil,err }
					method := proxyMethod{name:member.Names[0].Name,signature:strings.TrimPrefix(signature.String(),"func"),
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
						ep,err := endpointArguments(set,fn,args,line)
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

func requestLine(doc *ast.CommentGroup) (endpointMethod,bool,error) {
	if doc==nil { return endpointMethod{},false,nil }
	found:=false
	var result endpointMethod
	for _,comment := range doc.List {
		line := strings.TrimSpace(strings.TrimPrefix(comment.Text,"//"))
		if !strings.HasPrefix(line,"@RequestLine") { continue }
		if line!="@RequestLine" && !strings.HasPrefix(line,"@RequestLine ") { continue }
		if found { return result,true,fmt.Errorf("duplicate @RequestLine") }
		found=true
		args := strings.Fields(strings.TrimSpace(strings.TrimPrefix(line,"@RequestLine")))
		if len(args)!=2 { return result,true,fmt.Errorf("@RequestLine requires METHOD /absolute/path") }
		result.verb,result.path=args[0],args[1]
		switch result.verb {
		case "GET","POST","PUT","PATCH","DELETE":
		default:
			return result,true,fmt.Errorf("unsupported @RequestLine verb %q",result.verb)
		}
		if !strings.HasPrefix(result.path,"/") || strings.Contains(result.path,"//") ||
			strings.ContainsAny(result.path,"?#") {
			return result,true,fmt.Errorf("invalid @RequestLine path %q",result.path)
		}
		for _,segment := range strings.Split(result.path,"/") {
			if segment=="."||segment==".." { return result,true,fmt.Errorf("unsafe @RequestLine path") }
		}
	}
	return result,found,nil
}

func endpointArguments(set *token.FileSet, fn *ast.FuncType, args []string, method endpointMethod) (endpointMethod,error) {
	if len(args)!=2 || fn.Params==nil || len(fn.Params.List)!=2 || fn.Results==nil {
		return method,fmt.Errorf("@RequestLine requires context and one operation-specific request struct")
	}
	if !selectorType(fn.Params.List[0].Type,"context","Context") {
		return method,fmt.Errorf("first @RequestLine parameter must be context.Context")
	}
	// A request DTO is a named type (local or qualified). Its only permitted
	// exported fields are Params, Query and RequestBody; the transport enforces
	// their shape before dispatch, without any generic Empty placeholders.
	switch fn.Params.List[1].Type.(type) {
	case *ast.SelectorExpr, *ast.Ident:
	default:
		return method,fmt.Errorf("second @RequestLine parameter must be a named request struct")
	}
	results:=fn.Results.List
	if len(results)!=2 && len(results)!=3 {return method,fmt.Errorf("expected (response,error) or (response,bool,error)")}
	if !identType(results[len(results)-1].Type,"error") {return method,fmt.Errorf("last result must be error")}
	method.created=len(results)==3
	if method.created && !identType(results[1].Type,"bool") {
		return method,fmt.Errorf("second result must be bool (created flag)")
	}
	if method.created && method.verb!="POST" {return method,fmt.Errorf("created flag is supported only for POST")}
	if len(results[0].Names)>0 || len(results[len(results)-1].Names)>0 {
		return method,fmt.Errorf("named @RequestLine results are not supported")
	}
	var result bytes.Buffer
	if err:=format.Node(&result,set,results[0].Type);err!=nil{return method,err}
	method.resultType=result.String()
	return method,nil
}

func selectorType(expr ast.Expr, pkg,name string) bool {
	selector,ok:=expr.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name!=name{return false}
	base,ok:=selector.X.(*ast.Ident)
	return ok && base.Name==pkg
}
func identType(expr ast.Expr,name string)bool{
	id,ok:=expr.(*ast.Ident)
	return ok && id.Name==name
}

func directive(doc *ast.CommentGroup,name string)(string,bool,error) {
	if doc==nil {return "",false,nil}
	key:="@"+name
	found:=false
	var value string
	for _,comment:=range doc.List {
		line:=strings.TrimSpace(strings.TrimPrefix(comment.Text,"//"))
		if line!=key&&!strings.HasPrefix(line,key+" ") {continue}
		if found {return "",true,fmt.Errorf("duplicate %s",key)}
		found=true
		value=strings.TrimSpace(strings.TrimPrefix(line,key))
		if value==""||strings.ContainsAny(value," \t\r\n"){return "",true,fmt.Errorf("%s expects one binding",key)}
	}
	return value,found,nil
}

func writeProxies(proxies []clientProxy,output string) error {
	if len(proxies)!=1{return fmt.Errorf("one @Client interface per service is supported")}
	proxy:=proxies[0]
	usesEndpoints:=false
	for _,m:=range proxy.methods { if m.endpoint!=nil {usesEndpoints=true} }
	var b strings.Builder
	b.WriteString("// Code generated by GrooveHub @Client/@Proxy/@RequestLine. DO NOT EDIT.\npackage impl\n\nimport (\n")
	for _,imp:=range proxy.imports {
		if imp.Name!=nil {fmt.Fprintf(&b,"%s ",imp.Name.Name)}
		fmt.Fprintf(&b,"%s\n",imp.Path.Value)
	}
	if usesEndpoints {b.WriteString("httpendpoint \"grovehub/libs/core/http/endpoint\"\n")}
	b.WriteString(")\n\n")
	for _,m:=range proxy.methods {
		fmt.Fprintf(&b,"func (c *%s) %s%s {\n",proxy.receiver,m.name,m.signature)
		if m.endpoint==nil {
			if m.hasResults {b.WriteString("return ")}
			fmt.Fprintf(&b,"c.%s(%s)\n}\n\n",m.delegate,m.arguments)
			continue
		}
		ep:=m.endpoint
		if ep.created {
			fmt.Fprintf(&b,"result, status, err := httpendpoint.CallStruct[%s](%s, c.handler, %q, %q, %s)\n",
				ep.resultType,strings.Split(m.arguments,", ")[0],ep.verb,ep.path,strings.Split(m.arguments,", ")[1])
			b.WriteString("return result, status == 201 && err == nil, err\n}\n\n")
		} else {
			fmt.Fprintf(&b,"result, _, err := httpendpoint.CallStruct[%s](%s, c.handler, %q, %q, %s)\n",
				ep.resultType,strings.Split(m.arguments,", ")[0],ep.verb,ep.path,strings.Split(m.arguments,", ")[1])
			b.WriteString("return result, err\n}\n\n")
		}
	}
	source,err:=format.Source([]byte(b.String()))
	if err!=nil {return fmt.Errorf("format generated client: %w",err)}
	previous,err:=os.ReadFile(output)
	if err==nil && bytes.Equal(previous,source){return nil}
	if err!=nil&&!os.IsNotExist(err){return err}
	if err:=os.MkdirAll(filepath.Dir(output),0755);err!=nil{return err}
	return os.WriteFile(output,source,0644)
}

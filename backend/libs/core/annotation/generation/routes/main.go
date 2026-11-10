// Command routes generates net/http registration from annotated controller interfaces.
// The generated file is stored under service-local .generated/routes and injected
// into the original routes package through the Go overlay.
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
	method, verb, path string
}
type routeGroup struct {
	name string
	routes []route
}

var verbs=map[string]bool{"GET":true,"POST":true,"PUT":true,"PATCH":true,"DELETE":true}

func main() {
	if len(os.Args)!=2 { fmt.Fprintln(os.Stderr,"usage: routes <backend-dir>"); os.Exit(2) }
	if err:=generateAll(os.Args[1]);err!=nil { fmt.Fprintln(os.Stderr,err);os.Exit(1) }
}

func generateAll(backend string) error {
	services,err:=os.ReadDir(filepath.Join(backend,"services"));if err!=nil {return err}
	for _,service:=range services {
		if !service.IsDir(){continue}
		name:=service.Name()
		directory:=filepath.Join(backend,"services",name,"controller")
		if info,err:=os.Stat(directory);err!=nil||!info.IsDir(){continue}
		groups,err:=discover(directory);if err!=nil{return fmt.Errorf("%s routes: %w",name,err)}
		output:=filepath.Join(backend,"services",name,".generated","routes","zz_routes.gen.go")
		if len(groups)==0 {
			if err:=os.Remove(output);err!=nil&&!os.IsNotExist(err){return err}
			continue
		}
		if err:=writeRoutes(name,groups,output);err!=nil{return err}
	}
	return nil
}

func discover(dir string)([]routeGroup,error) {
	files,err:=filepath.Glob(filepath.Join(dir,"*.go"));if err!=nil{return nil,err}
	sort.Strings(files)
	var result []routeGroup
	seen:=map[string]string{}
	for _,path:=range files {
		if strings.HasSuffix(path,"_test.go"){continue}
		astFile,err:=parser.ParseFile(token.NewFileSet(),path,nil,parser.ParseComments);if err!=nil{return nil,err}
		for _,decl:=range astFile.Decls {
			gen,ok:=decl.(*ast.GenDecl);if !ok||gen.Tok!=token.TYPE{continue}
			for _,node:=range gen.Specs {
				spec,ok:=node.(*ast.TypeSpec);if !ok {continue}
				iface,ok:=spec.Type.(*ast.InterfaceType);if !ok {continue}
				doc:=spec.Doc;if doc==nil{doc=gen.Doc}
				prefix,found,err:=prefixFor(doc);if err!=nil{return nil,fmt.Errorf("%s: %w",spec.Name.Name,err)}
				if !found{continue}
				if !strings.HasSuffix(spec.Name.Name,"Controller"){return nil,fmt.Errorf("@Routes interface %s must end with Controller",spec.Name.Name)}
				group:=routeGroup{name:spec.Name.Name}
				for _,member:=range iface.Methods.List {
					if len(member.Names)!=1{return nil,fmt.Errorf("%s cannot embed interfaces",spec.Name.Name)}
					if !validSignature(member.Type){return nil,fmt.Errorf("%s.%s must accept http.ResponseWriter and *http.Request",spec.Name.Name,member.Names[0].Name)}
					r,err:=methodFor(member.Doc);if err!=nil{return nil,err}
					if r.verb==""{return nil,fmt.Errorf("%s.%s missing @GET/@POST/etc",spec.Name.Name,member.Names[0].Name)}
					r.method=member.Names[0].Name
					if !strings.HasPrefix(r.path,"/"){return nil,fmt.Errorf("%s.%s path must begin with /",spec.Name.Name,r.method)}
					if r.path=="/"{r.path=prefix}else{r.path=strings.TrimRight(prefix,"/")+r.path}
					key:=r.verb+" "+r.path
					if owner,ok:=seen[key];ok{return nil,fmt.Errorf("duplicate route %s: %s and %s.%s",key,owner,spec.Name.Name,r.method)}
					seen[key]=spec.Name.Name+"."+r.method
					group.routes=append(group.routes,r)
				}
				result=append(result,group)
			}
		}
	}
	sort.Slice(result,func(i,j int)bool{return result[i].name<result[j].name})
	return result,nil
}

func prefixFor(doc *ast.CommentGroup)(string,bool,error) {
	if doc==nil{return "",false,nil}
	var prefix string
	found:=false
	for _,comment:=range doc.List {
		line:=strings.TrimSpace(strings.TrimPrefix(comment.Text,"//"))
		if line=="@Routes"||strings.HasPrefix(line,"@Routes "){
			if found{return "",true,fmt.Errorf("duplicate @Routes declaration")}
			found=true
			prefix=strings.TrimSpace(strings.TrimPrefix(line,"@Routes"))
			if !strings.HasPrefix(prefix,"/")||strings.ContainsAny(prefix," \t\r\n")||strings.Contains(prefix,"//"){return "",true,fmt.Errorf("invalid route prefix %q",prefix)}
		}
	}
	return prefix,found,nil
}

func methodFor(doc *ast.CommentGroup)(route,error) {
	var r route
	if doc==nil{return r,nil}
	for _,comment:=range doc.List {
		line:=strings.TrimSpace(strings.TrimPrefix(comment.Text,"//"))
		if !strings.HasPrefix(line,"@"){continue}
		parts:=strings.Fields(strings.TrimPrefix(line,"@"))
		if len(parts)==0||!verbs[parts[0]]{continue}
		if r.verb!="" {return r,fmt.Errorf("multiple HTTP verb annotations")}
		if len(parts)!=2{return r,fmt.Errorf("@%s requires exactly one path",parts[0])}
		r.verb=parts[0];r.path=parts[1]
	}
	return r,nil
}

func validSignature(expr ast.Expr)bool {
	fn,ok:=expr.(*ast.FuncType);if !ok||fn.Results!=nil||fn.Params==nil||len(fn.Params.List)!=2{return false}
	writer,ok:=fn.Params.List[0].Type.(*ast.SelectorExpr);if !ok||writer.Sel.Name!="ResponseWriter"{return false}
	pkg,ok:=writer.X.(*ast.Ident);if !ok||pkg.Name!="http"{return false}
	pointer,ok:=fn.Params.List[1].Type.(*ast.StarExpr);if !ok{return false}
	request,ok:=pointer.X.(*ast.SelectorExpr);if !ok||request.Sel.Name!="Request"{return false}
	pkg,ok=request.X.(*ast.Ident);return ok&&pkg.Name=="http"
}

func writeRoutes(service string, groups []routeGroup, output string)error {
	var b strings.Builder
	b.WriteString("// Code generated by GrooveHub @Routes. DO NOT EDIT.\npackage routes\n\n")
	fmt.Fprintf(&b,"import (\n\"net/http\"\nhttptracing \"grovehub/libs/core/http/tracing\"\n\"grovehub/services/%s/controller\"\n)\n\n",service)
	for _,group:=range groups {
		fmt.Fprintf(&b,"func register%sRoutes(mux *http.ServeMux, c controller.%s) {\n",strings.TrimSuffix(group.name,"Controller"),group.name)
		for _,r:=range group.routes {
			pattern:=r.verb+" "+r.path
			fmt.Fprintf(&b,"mux.Handle(%q, httptracing.ServiceHandler(%q,%q,http.HandlerFunc(c.%s)))\n",pattern,service,pattern,r.method)
		}
		b.WriteString("}\n\n")
	}
	formatted,err:=format.Source([]byte(b.String()));if err!=nil{return err}
	previous,err:=os.ReadFile(output)
	if err==nil&&bytes.Equal(previous,formatted){return nil}
	if err!=nil&&!os.IsNotExist(err){return err}
	if err:=os.MkdirAll(filepath.Dir(output),0755);err!=nil{return err}
	return os.WriteFile(output,formatted,0644)
}

// Command search generates repository search plumbing from an explicit
// @Searchable declaration. SQL field expressions stay developer-owned in
// coresearch.Config; no query, join, or field is inferred from user input.
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

type spec struct {
	service, impl, entity, config, alias string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: search <backend-dir>")
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
	for _, service := range services {
		if !service.IsDir() { continue }
		owner := filepath.Join(backend, "services", service.Name())
		dir := filepath.Join(owner, "repository")
		if _, err := os.Stat(dir); os.IsNotExist(err) { continue } else if err != nil { return err }
		specs, err := discover(dir, service.Name())
		if err != nil { return fmt.Errorf("%s search: %w", service.Name(), err) }
		out := filepath.Join(owner, ".generated", "repository", "impl", "zz_search.gen.go")
		if len(specs)==0 {
			if err:=os.Remove(out);err!=nil&&!os.IsNotExist(err) {return err}
			continue
		}
		if err := validateImplementation(filepath.Join(dir,"impl"), specs); err != nil { return err }
		if err := write(out, specs); err != nil { return err }
	}
	return nil
}

func discover(dir, service string) ([]spec, error) {
	paths, err := filepath.Glob(filepath.Join(dir,"*.go"))
	if err != nil { return nil, err }
	sort.Strings(paths)
	var result []spec
	for _, path := range paths {
		if strings.HasSuffix(path,"_test.go") { continue }
		file, err := parser.ParseFile(token.NewFileSet(),path,nil,parser.ParseComments)
		if err != nil { return nil, err }
		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok!=token.TYPE { continue }
			for _, item := range gen.Specs {
				typeSpec, ok := item.(*ast.TypeSpec)
				if !ok { continue }
				iface, ok := typeSpec.Type.(*ast.InterfaceType)
				if !ok { continue }
				doc := typeSpec.Doc
				if doc==nil {doc=gen.Doc}
				repoDirective := directive(doc,"Repository")
				for _, member := range iface.Methods.List {
					args := directive(member.Doc,"Searchable")
					if args=="" {continue}
					opts,err:=options(args,[]string{"impl","entity","config","alias"})
					if err!=nil{return nil,fmt.Errorf("%s.%s: %w",typeSpec.Name.Name,member.Names[0].Name,err)}
					if len(member.Names)!=1 || member.Names[0].Name!="Search" {
						return nil,fmt.Errorf("%s: @Searchable only supports a Search method",typeSpec.Name.Name)
					}
					if opts["entity"]==""||opts["impl"]==""||opts["config"]==""||opts["alias"]=="" { return nil,fmt.Errorf("%s: incomplete search config",typeSpec.Name.Name) }
					expectedRepo:="entity="+opts["entity"]
					if strings.TrimSpace(repoDirective)!=expectedRepo {return nil,fmt.Errorf("%s: @Repository %s must match @Searchable entity",typeSpec.Name.Name,expectedRepo)}
					if opts["impl"]!=typeSpec.Name.Name {return nil,fmt.Errorf("%s.Search: impl must equal repository interface name",typeSpec.Name.Name)}
					if err:=validateSignature(member.Type,opts["entity"]);err!=nil{return nil,fmt.Errorf("%s.Search: %w",typeSpec.Name.Name,err)}
					filterFields,err:=fieldsFrom(member.Doc,"Filterable")
					if err!=nil{return nil,err}
					sortFields,err:=fieldsFrom(member.Doc,"Sortable")
					if err!=nil{return nil,err}
					cursorArgs,err:=options(directive(member.Doc,"Cursor"),[]string{"field"})
					if err!=nil{return nil,fmt.Errorf("@Cursor: %w",err)}
					if !contains(sortFields,cursorArgs["field"]) {return nil,fmt.Errorf("@Cursor field must be sortable")}
					implDir:=filepath.Join(dir,"impl")
					if err:=validateConfig(implDir,opts["config"],filterFields,sortFields,cursorArgs["field"]);err!=nil{return nil,err}
					result=append(result,spec{service:service,impl:opts["impl"],entity:opts["entity"],config:opts["config"],alias:opts["alias"]})
				}
			}
		}
	}
	return result,nil
}

func directive(group *ast.CommentGroup, name string) string {
	if group==nil {return ""}
	for _, c:=range group.List {
		line:=strings.TrimSpace(strings.TrimPrefix(c.Text,"//"))
		if line=="@"+name {return ""}
		if strings.HasPrefix(line,"@"+name+" ") {return strings.TrimSpace(strings.TrimPrefix(line,"@"+name))}
	}
	return ""
}

func options(raw string, expected []string)(map[string]string,error) {
	r:=make(map[string]string,len(expected))
	for _,part:=range strings.Fields(raw) {
		p:=strings.SplitN(part,"=",2)
		if len(p)!=2||p[1]=="" {return nil,fmt.Errorf("invalid annotation argument %q",part)}
		valid:=false
		for _,key:=range expected {if p[0]==key {valid=true}}
		if !valid||r[p[0]]!="" {return nil,fmt.Errorf("unknown or repeated annotation option %s",p[0])}
		r[p[0]]=p[1]
	}
	for _,key:=range expected {if r[key]=="" {return nil,fmt.Errorf("missing %s option",key)}}
	return r,nil
}

func fieldsFrom(group *ast.CommentGroup, name string) ([]string,error) {
	args,err:=options(directive(group,name),[]string{"fields"})
	if err!=nil {return nil,fmt.Errorf("@%s: %w",name,err)}
	return strings.Split(args["fields"],","),nil
}

func validateSignature(expr ast.Expr,entity string)error {
	fn,ok:=expr.(*ast.FuncType)
	if !ok||fn.Params==nil||fn.Results==nil||len(fn.Params.List)!=2||len(fn.Results.List)!=2 {
		return fmt.Errorf("requires (context.Context, coresearch.Query[string]) (coresearch.QueryResult[entity.%s], error)",entity)
	}
	if exprString(fn.Params.List[0].Type)!="context.Context" ||
		exprString(fn.Params.List[1].Type)!="coresearch.Query[string]" ||
		exprString(fn.Results.List[0].Type)!="coresearch.QueryResult[entity."+entity+"]" ||
		exprString(fn.Results.List[1].Type)!="error" {
		return fmt.Errorf("unsupported signature: expected context.Context, coresearch.Query[string] and QueryResult[entity.%s]",entity)
	}
	return nil
}
func exprString(e ast.Expr) string {
	var b bytes.Buffer
	if err:=format.Node(&b,token.NewFileSet(),e);err!=nil{return ""}
	return strings.Join(strings.Fields(b.String()),"")
}

func findConfig(dir,name string) (*ast.CompositeLit,error) {
	paths,err:=filepath.Glob(filepath.Join(dir,"*.go"))
	if err!=nil{return nil,err}
	for _,p:=range paths {
		file,err:=parser.ParseFile(token.NewFileSet(),p,nil,0)
		if err!=nil{return nil,err}
		for _,dec:=range file.Decls {
			g,ok:=dec.(*ast.GenDecl);if !ok||g.Tok!=token.VAR{continue}
			for _,item:=range g.Specs {
				value,ok:=item.(*ast.ValueSpec);if !ok||len(value.Values)!=1{continue}
				for _,ident:=range value.Names {
					if ident.Name!=name{continue}
					lit,ok:=value.Values[0].(*ast.CompositeLit)
					if !ok||exprString(lit.Type)!="coresearch.Config" {return nil,fmt.Errorf("%s must be coresearch.Config literal",name)}
					return lit,nil
				}
			}
		}
	}
	return nil,fmt.Errorf("search configuration %s not found",name)
}
func field(lit *ast.CompositeLit,key string)ast.Expr{
	for _,item:=range lit.Elts {
		pair,ok:=item.(*ast.KeyValueExpr);if !ok {continue}
		id,ok:=pair.Key.(*ast.Ident);if ok&&id.Name==key{return pair.Value}
	}
	return nil
}
func stringExpr(node ast.Expr)(string,error){
	value,ok:=node.(*ast.BasicLit);if !ok||value.Kind!=token.STRING{return "",fmt.Errorf("expected constant string expression")}
	return strconv.Unquote(value.Value)
}
func names(node ast.Expr)([]string,error){
	lit,ok:=node.(*ast.CompositeLit);if !ok{return nil,fmt.Errorf("expected static map literal")}
	names:=make([]string,0,len(lit.Elts))
	for _,item:=range lit.Elts {
		pair,ok:=item.(*ast.KeyValueExpr);if !ok{return nil,fmt.Errorf("map requires keyed entries")}
		s,err:=stringExpr(pair.Key);if err!=nil{return nil,err}
		names=append(names,s)
	}
	sort.Strings(names)
	return names,nil
}
func contains(all []string,target string)bool{
	for _,v:=range all {if v==target{return true}}
	return false
}
func compareFields(actual,expected []string,label string)error{
	sort.Strings(actual)
	sort.Strings(expected)
	if strings.Join(actual,",")!=strings.Join(expected,","){
		return fmt.Errorf("%s annotation must exactly match configured allowlist: configured=%v, declared=%v",label,actual,expected)
	}
	return nil
}
func validateConfig(dir,config string,filters,sorts []string,cursor string) error {
	lit,err:=findConfig(dir,config);if err!=nil{return err}
	scalars,err:=names(field(lit,"ScalarFields"));if err!=nil{return err}
	relations,err:=names(field(lit,"RelationFields"));if err!=nil{return err}
	if err:=compareFields(append(scalars,relations...),filters,"@Filterable");err!=nil{return err}
	configuredSorts,err:=names(field(lit,"SortFields"));if err!=nil{return err}
	if err:=compareFields(configuredSorts,sorts,"@Sortable");err!=nil{return err}
	cursorExpr,err:=stringExpr(field(lit,"CursorExpression"));if err!=nil{return err}
	sortsNode,ok:=field(lit,"SortFields").(*ast.CompositeLit);if !ok{return fmt.Errorf("invalid SortFields")}
	for _,item:=range sortsNode.Elts {
		pair:=item.(*ast.KeyValueExpr)
		name,err:=stringExpr(pair.Key);if err!=nil{return err}
		if name!=cursor {continue}
		value,ok:=pair.Value.(*ast.CompositeLit);if !ok{return fmt.Errorf("sort field %s must use a literal",cursor)}
		expression,err:=stringExpr(field(value,"Expression"));if err!=nil{return err}
		if expression!=cursorExpr{return fmt.Errorf("@Cursor field %s must match CursorExpression",cursor)}
		return nil
	}
	return fmt.Errorf("@Cursor sort %s missing",cursor)
}
func validateImplementation(dir string,specs []spec)error{
	files,err:=filepath.Glob(filepath.Join(dir,"*.go"));if err!=nil{return err}
	structs:=make(map[string]bool)
	for _,p:=range files {
		f,err:=parser.ParseFile(token.NewFileSet(),p,nil,0);if err!=nil{return err}
		for _,d:=range f.Decls {
			switch n:=d.(type) {
			case *ast.GenDecl:
				if n.Tok!=token.TYPE{continue}
				for _,item:=range n.Specs {
					ts,ok:=item.(*ast.TypeSpec);if !ok{continue}
					st,ok:=ts.Type.(*ast.StructType);if !ok{continue}
					for _,v:=range st.Fields.List {for _,id:=range v.Names {if id.Name=="db"{structs[ts.Name.Name]=true}}}
				}
			case *ast.FuncDecl:
				if n.Recv==nil||n.Name.Name!="Search"{continue}
				for _,wanted:=range specs {
					if strings.TrimPrefix(strings.TrimPrefix(exprString(n.Recv.List[0].Type),"*"),"impl.")==wanted.impl {
						return fmt.Errorf("%s has a manual Search implementation; remove it when using @Searchable",wanted.impl)
					}
				}
			}
		}
	}
	for _,wanted:=range specs {if !structs[wanted.impl]{return fmt.Errorf("%s needs a db field on its implementation struct",wanted.impl)}}
	return nil
}
func write(path string,specs []spec)error{
	var b strings.Builder
	b.WriteString("// Code generated by GrooveHub @Searchable. DO NOT EDIT.\npackage impl\n\n")
	b.WriteString("import (\n\"context\"\ncoresearch \"grovehub/libs/core/search\"\ncoretransaction \"grovehub/libs/core/database/transaction\"\n")
	// All @Searchable declarations in this service share its entity package.
	fmt.Fprintf(&b,"\"grovehub/services/%s/entity\"\n)\n\n",specs[0].service)
	for _,sp:=range specs{
		if sp.service!=specs[0].service{return fmt.Errorf("mixed service annotations")}
		fmt.Fprintf(&b,"func (repository *%s) Search(ctx context.Context, query coresearch.Query[string]) (coresearch.QueryResult[entity.%s], error) {\n",sp.impl,sp.entity)
		fmt.Fprintf(&b,"return coresearch.Execute[string, entity.%s](coretransaction.Database(ctx, repository.db).Model(&entity.%s{}).Table(entity.%s{}.TableName()+\" AS %s\").Select(\"%s.*\"), query, %s)\n}\n\n",sp.entity,sp.entity,sp.entity,sp.alias,sp.alias,sp.config)
	}
	formatted,err:=format.Source([]byte(b.String()))
	if err!=nil{return fmt.Errorf("format generated search code: %w",err)}
	if err:=os.MkdirAll(filepath.Dir(path),0755);err!=nil{return err}
	if old,err:=os.ReadFile(path);err==nil&&bytes.Equal(old,formatted){return nil}
	return os.WriteFile(path,formatted,0644)
}

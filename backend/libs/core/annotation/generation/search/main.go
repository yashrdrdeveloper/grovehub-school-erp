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
	configImport, configAlias string
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
					opts,err:=searchOptions(args)
					if err!=nil{return nil,fmt.Errorf("%s.%s: %w",typeSpec.Name.Name,member.Names[0].Name,err)}
					if len(member.Names)!=1 || member.Names[0].Name!="Search" {
						return nil,fmt.Errorf("%s: @Searchable only supports a Search method",typeSpec.Name.Name)
					}
					repoOpts,err:=options(repoDirective,[]string{"entity"})
					if err!=nil{return nil,fmt.Errorf("%s.Search: @Repository entity declaration required: %w",typeSpec.Name.Name,err)}
					entityName:=repoOpts["entity"]
					if opts["entity"]!="" && opts["entity"]!=entityName {
						return nil,fmt.Errorf("%s.Search: @Searchable entity must match @Repository",typeSpec.Name.Name)
					}
					implName:=typeSpec.Name.Name
					if opts["impl"]!="" && opts["impl"]!=implName {
						return nil,fmt.Errorf("%s.Search: @Searchable impl must match repository interface",typeSpec.Name.Name)
					}
					if err:=validateSignature(member.Type,entityName);err!=nil{return nil,fmt.Errorf("%s.Search: %w",typeSpec.Name.Name,err)}
					filterFields,err:=optionalFields(member.Doc,"Filterable")
					if err!=nil{return nil,err}
					sortFields,err:=optionalFields(member.Doc,"Sortable")
					if err!=nil{return nil,err}
					cursorName:=""
					if raw:=directive(member.Doc,"Cursor");raw!="" {
						cursorArgs,err:=options(raw,[]string{"field"})
						if err!=nil{return nil,fmt.Errorf("@Cursor: %w",err)}
						cursorName=cursorArgs["field"]
						if len(sortFields)>0 && !contains(sortFields,cursorName) {return nil,fmt.Errorf("@Cursor field must be sortable")}
					}
					configDir,configImport,configAlias,configName,err:=resolveConfig(dir,service,opts["config"])
					if err!=nil{return nil,err}
					if err:=validateConfig(configDir,configName,filterFields,sortFields,cursorName);err!=nil{return nil,err}
					result=append(result,spec{
						service:service,impl:implName,entity:entityName,config:opts["config"],
						alias:opts["alias"],configImport:configImport,configAlias:configAlias,
					})
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

func searchOptions(raw string) (map[string]string,error) {
	result:=map[string]string{}
	allowed:=map[string]bool{"impl":true,"entity":true,"config":true,"alias":true}
	for _,item:=range strings.Fields(raw) {
		pair:=strings.SplitN(item,"=",2)
		if len(pair)!=2||pair[1]==""||!allowed[pair[0]]||result[pair[0]]!="" {
			return nil,fmt.Errorf("invalid or repeated @Searchable option %q",item)
		}
		result[pair[0]]=pair[1]
	}
	if result["config"]==""||result["alias"]=="" {
		return nil,fmt.Errorf("@Searchable requires config=... alias=...")
	}
	return result,nil
}

func optionalFields(group *ast.CommentGroup,name string)([]string,error){
	raw:=directive(group,name)
	if raw=="" {return nil,nil}
	opts,err:=options(raw,[]string{"fields"})
	if err!=nil{return nil,fmt.Errorf("@%s: %w",name,err)}
	return strings.Split(opts["fields"],","),nil
}

// resolveConfig finds an explicitly imported repository search configuration
// in a subfolder by its declared Go package, not by an inferred SQL/table name.
func resolveConfig(repositoryDir,service,reference string)(dir,importPath,alias,variable string,err error) {
	parts:=strings.Split(reference,".")
	if len(parts)==1 {
		if !token.IsIdentifier(parts[0]) {return "","","","",fmt.Errorf("invalid @Searchable config %q",reference)}
		return filepath.Join(repositoryDir,"impl"),"","",parts[0],nil
	}
	if len(parts)!=2||!token.IsIdentifier(parts[0])||!token.IsIdentifier(parts[1]) {
		return "","","","",fmt.Errorf("invalid @Searchable config %q",reference)
	}
	entries,err:=os.ReadDir(repositoryDir)
	if err!=nil{return "","","","",err}
	for _,entry:=range entries {
		if !entry.IsDir() {continue}
		candidate:=filepath.Join(repositoryDir,entry.Name())
		files,err:=filepath.Glob(filepath.Join(candidate,"*.go"))
		if err!=nil{return "","","","",err}
		for _,f:=range files {
			pkg,err:=parser.ParseFile(token.NewFileSet(),f,nil,parser.PackageClauseOnly)
			if err!=nil{return "","","","",err}
			if pkg.Name.Name!=parts[0] {continue}
			return candidate,"grovehub/services/"+service+"/repository/"+entry.Name(),parts[0],parts[1],nil
		}
	}
	return "","","","",fmt.Errorf("@Searchable config package %q not found under %s",parts[0],repositoryDir)
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
	if len(filters)>0 {
		if err:=compareFields(append(scalars,relations...),filters,"@Filterable");err!=nil{return err}
	}
	configuredSorts,err:=names(field(lit,"SortFields"));if err!=nil{return err}
	if len(sorts)>0 {
		if err:=compareFields(configuredSorts,sorts,"@Sortable");err!=nil{return err}
	}
	cursorExpr,err:=stringExpr(field(lit,"CursorExpression"));if err!=nil{return err}
	sortsNode,ok:=field(lit,"SortFields").(*ast.CompositeLit);if !ok{return fmt.Errorf("invalid SortFields")}
	for _,item:=range sortsNode.Elts {
		pair:=item.(*ast.KeyValueExpr)
		name,err:=stringExpr(pair.Key);if err!=nil{return err}
		if cursor!="" && name!=cursor {continue}
		value,ok:=pair.Value.(*ast.CompositeLit);if !ok{return fmt.Errorf("sort field %s must use a literal",cursor)}
		expression,err:=stringExpr(field(value,"Expression"));if err!=nil{return err}
		if expression==cursorExpr{return nil}
	}
	if cursor!="" {
		return fmt.Errorf("@Cursor sort %s must match CursorExpression",cursor)
	}
	return fmt.Errorf("no sortable field matches the stable CursorExpression %q",cursorExpr)
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
	fmt.Fprintf(&b,"\"grovehub/services/%s/entity\"\n",specs[0].service)
	imports:=map[string]string{}
	for _,sp:=range specs {
		if sp.configImport=="" {continue}
		if old,ok:=imports[sp.configAlias];ok&&old!=sp.configImport {
			return fmt.Errorf("conflicting config package alias %s",sp.configAlias)
		}
		imports[sp.configAlias]=sp.configImport
	}
	aliases:=make([]string,0,len(imports))
	for name:=range imports {aliases=append(aliases,name)}
	sort.Strings(aliases)
	for _,name:=range aliases {fmt.Fprintf(&b,"%s %q\n",name,imports[name])}
	b.WriteString(")\n\n")
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

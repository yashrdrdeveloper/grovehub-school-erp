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
        if info, err := os.Stat(filepath.Join(owner, "internal", "app", "repository")); err == nil && info.IsDir() {
            dir = filepath.Join(owner, "internal", "app", "repository")
        }
		if _, err := os.Stat(dir); os.IsNotExist(err) { continue } else if err != nil { return err }
		// Search contracts are scoped to repository/<resource>, preserving
		// repository/impl as a legacy flat layout for other services.
		entries,err:=os.ReadDir(dir)
		if err!=nil{return err}
		directories:=[]string{dir}
		for _,entry:=range entries {
			if !entry.IsDir(){continue}
			if entry.Name()=="impl"||entry.Name()=="constants"||entry.Name()=="helpers"||strings.HasPrefix(entry.Name(),"."){continue}
			directories=append(directories,filepath.Join(dir,entry.Name()))
		}
		for _,resourceDir:=range directories {
			specs,err:=discover(resourceDir,service.Name())
			if err!=nil{return fmt.Errorf("%s search %s: %w",service.Name(),resourceDir,err)}
			rel,err:=filepath.Rel(dir,resourceDir)
			if err!=nil{return err}
			rootRel, err := filepath.Rel(owner, dir)
            if err != nil { return err }
            outDir:=filepath.Join(owner,".generated",rootRel)
			if rel!="." {outDir=filepath.Join(outDir,rel)}
			out:=filepath.Join(outDir,"impl","zz_search.gen.go")
			if len(specs)==0 {
				if err:=os.Remove(out);err!=nil&&!os.IsNotExist(err){return err}
				continue
			}
			if err:=validateImplementation(filepath.Join(resourceDir,"impl"),specs);err!=nil{return err}
			if err:=write(out,specs);err!=nil{return err}
		}
	}
	return nil
}

// @Searchable belongs to the concrete Search method, not the public interface.
// The developer-maintained Search method delegates to generated searchGenerated.
func discover(dir, service string) ([]spec, error) {
	if err := rejectInterfaceSearchables(dir); err != nil { return nil, err }
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
			entityName, err := repositoryEntity(dir, receiverType)
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
			configDir, configImport, configAlias, configName, err := resolveConfig(dir, service, opts["config"])
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
func rejectInterfaceSearchables(dir string) error {
	paths, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil { return err }
	for _, filename := range paths {
		if strings.HasSuffix(filename, "_test.go") { continue }
		file, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
		if err != nil { return err }
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
						return fmt.Errorf("%s: move @Searchable from %s interface to the concrete Search method", filename, t.Name.Name)
					}
				}
			}
		}
	}
	return nil
}

// repositoryEntity ensures the annotated implementation satisfies the
// explicit repository interface and its @Repository ownership declaration.
func repositoryEntity(dir, receiver string) (string, error) {
	paths, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil { return "", err }
	for _, filename := range paths {
		if strings.HasSuffix(filename, "_test.go") { continue }
		file, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
		if err != nil { return "", err }
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

func hasDirective(doc *ast.CommentGroup, name string) bool {
	if doc == nil { return false }
	for _, comment := range doc.List {
		text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
		if text == "@"+name || strings.HasPrefix(text, "@"+name+" ") { return true }
	}
	return false
}

// The explicit Search method is intentionally small: config and query
// mechanics remain generated, and callers see an ordinary method in impl.
func validateSearchDelegate(fn *ast.FuncDecl) error {
	if fn.Body == nil || len(fn.Body.List) != 1 {
		return fmt.Errorf("@Searchable Search must delegate to repository.searchGenerated(ctx, query)")
	}
	stmt, ok := fn.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(stmt.Results) != 1 {
		return fmt.Errorf("@Searchable Search must return repository.searchGenerated(ctx, query)")
	}
	call, ok := stmt.Results[0].(*ast.CallExpr)
	if !ok || len(call.Args) != 2 {
		return fmt.Errorf("@Searchable Search must call the generated search helper")
	}
	target, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || target.Sel.Name != "searchGenerated" { return fmt.Errorf("@Searchable Search must use searchGenerated") }
	owner, ok := target.X.(*ast.Ident)
	if !ok || len(fn.Recv.List[0].Names) != 1 || owner.Name != fn.Recv.List[0].Names[0].Name {
		return fmt.Errorf("@Searchable Search must delegate using its repository receiver")
	}
	if exprString(call.Args[0]) != "ctx" || exprString(call.Args[1]) != "query" {
		return fmt.Errorf("@Searchable Search must forward ctx and query unchanged")
	}
	return nil
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
	// Repository-owned search configs may live in repository/search (legacy)
	// or repository/helpers/search, keeping the annotation independent of folders.
	candidates:=[]string{parts[0], filepath.Join("helpers",parts[0])}
	for _,relative:=range candidates {
		candidate:=filepath.Join(repositoryDir,relative)
		files,err:=filepath.Glob(filepath.Join(candidate,"*.go"))
		if err!=nil{return "","","","",err}
		for _,f:=range files {
			pkg,err:=parser.ParseFile(token.NewFileSet(),f,nil,parser.PackageClauseOnly)
			if err!=nil{return "","","","",err}
			if pkg.Name.Name!=parts[0] {continue}
			prefix:="repository/"
            if strings.Contains(filepath.ToSlash(repositoryDir),"/internal/app/repository"){prefix="internal/app/repository/"}
			if filepath.Base(repositoryDir)!="repository"{prefix+=filepath.Base(repositoryDir)+"/"}
			return candidate,"grovehub/services/"+service+"/"+prefix+filepath.ToSlash(relative),parts[0],parts[1],nil
		}
	}
	return "","","","",fmt.Errorf("@Searchable config package %q not found under %s or its helpers directory",parts[0],repositoryDir)
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
func validateImplementation(dir string, specs []spec) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil { return err }
	structs := make(map[string]bool)
	for _, filename := range files {
		file, err := parser.ParseFile(token.NewFileSet(), filename, nil, 0)
		if err != nil { return err }
		for _, decl := range file.Decls {
			switch node := decl.(type) {
			case *ast.GenDecl:
				if node.Tok != token.TYPE { continue }
				for _, item := range node.Specs {
					typeSpec, ok := item.(*ast.TypeSpec)
					if !ok { continue }
					st, ok := typeSpec.Type.(*ast.StructType)
					if !ok { continue }
					for _, field := range st.Fields.List {
						for _, name := range field.Names {
							if name.Name == "db" { structs[typeSpec.Name.Name] = true }
						}
					}
				}
			case *ast.FuncDecl:
				if node.Recv == nil || node.Name.Name != "searchGenerated" { continue }
				for _, wanted := range specs {
					if strings.TrimPrefix(exprString(node.Recv.List[0].Type), "*") == wanted.impl {
						return fmt.Errorf("%s already implements searchGenerated; @Searchable owns that helper", wanted.impl)
					}
				}
			}
		}
	}
	for _, wanted := range specs {
		if !structs[wanted.impl] { return fmt.Errorf("%s needs a db field on its implementation struct", wanted.impl) }
	}
	return nil
}

func write(path string,specs []spec)error{
	var b strings.Builder
	b.WriteString("// Code generated by GrooveHub @Searchable. DO NOT EDIT.\npackage impl\n\n")
	b.WriteString("import (\n\"context\"\ncoresearch \"grovehub/libs/core/search\"\ncoretransaction \"grovehub/libs/core/database/transaction\"\n")
	// All @Searchable declarations in this service share its entity package.
	entityPath := "grovehub/services/" + specs[0].service + "/entity"
    if strings.Contains(filepath.ToSlash(path), "/internal/app/repository/") {
        entityPath = "grovehub/services/" + specs[0].service + "/internal/app/entity"
    }
    fmt.Fprintf(&b,"%q\n",entityPath)
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
		fmt.Fprintf(&b,"func (repository *%s) searchGenerated(ctx context.Context, query coresearch.Query[string]) (coresearch.QueryResult[entity.%s], error) {\n",sp.impl,sp.entity)
		fmt.Fprintf(&b,"return coresearch.Execute[string, entity.%s](coretransaction.Database(ctx, repository.db).Model(&entity.%s{}).Table(entity.%s{}.TableName()+\" AS %s\").Select(\"%s.*\"), query, %s)\n}\n\n",sp.entity,sp.entity,sp.entity,sp.alias,sp.alias,sp.config)
	}
	formatted,err:=format.Source([]byte(b.String()))
	if err!=nil{return fmt.Errorf("format generated search code: %w",err)}
	if err:=os.MkdirAll(filepath.Dir(path),0755);err!=nil{return err}
	if old,err:=os.ReadFile(path);err==nil&&bytes.Equal(old,formatted){return nil}
	return os.WriteFile(path,formatted,0644)
}

// Command mappers generates strict typed DTO projection methods from @Mapper
// interfaces. Every destination field must have an explicit @Mapping rule.
// Business-specific transformations remain in maintained mapper methods.
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
	"sort"
	"strconv"
	"strings"
)

type mappingMethod struct {
	name, signature, output string
	fields []fieldMapping
}
type fieldMapping struct {name, expression string}
type mappingGroup struct {
	receiver string
	imports []*ast.ImportSpec
	methods []mappingMethod
}

func main() {
	if len(os.Args)!=2 {fmt.Fprintln(os.Stderr,"usage: mappers <backend-dir>");os.Exit(2)}
	if err:=generateAll(os.Args[1]);err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(1)}
}

func generateAll(backend string) error {
	services,err:=os.ReadDir(filepath.Join(backend,"services"));if err!=nil{return err}
	for _,entry:=range services {
		if !entry.IsDir(){continue}
		owner:=filepath.Join(backend,"services",entry.Name())
		dir:=filepath.Join(owner,"mappers")
		if stat,err:=os.Stat(dir);err!=nil||!stat.IsDir(){continue}
		// Generate once per resource package, retaining legacy flat mappers.
		entries,err:=os.ReadDir(dir)
		if err!=nil{return err}
		directories:=[]string{dir}
		for _,entry:=range entries {
			if entry.IsDir() && entry.Name()!="impl" && !strings.HasPrefix(entry.Name(),"."){
				directories=append(directories,filepath.Join(dir,entry.Name()))
			}
		}
		for _,resourceDir:=range directories {
			group,found,err:=discover(backend,resourceDir)
			if err!=nil{return fmt.Errorf("%s mapper %s: %w",entry.Name(),resourceDir,err)}
			rel,err:=filepath.Rel(dir,resourceDir)
			if err!=nil{return err}
			outDir:=filepath.Join(owner,".generated","mappers")
			if rel!="." {outDir=filepath.Join(outDir,rel)}
			out:=filepath.Join(outDir,"impl","zz_mappers.gen.go")
			if !found||len(group.methods)==0 {
				if err:=os.Remove(out);err!=nil&&!os.IsNotExist(err){return err}
				continue
			}
			if err:=write(group,out);err!=nil{return fmt.Errorf("%s mapper %s: %w",entry.Name(),resourceDir,err)}
		}
	}
	return nil
}

func discover(backend,dir string)(mappingGroup,bool,error) {
	files,err:=filepath.Glob(filepath.Join(dir,"*.go"));if err!=nil{return mappingGroup{},false,err}
	sort.Strings(files)
	var result mappingGroup
	found:=false
	for _,filePath:=range files {
		if strings.HasSuffix(filePath,"_test.go"){continue}
		set:=token.NewFileSet()
		file,err:=parser.ParseFile(set,filePath,nil,parser.ParseComments);if err!=nil{return result,false,err}
		for _,decl:=range file.Decls {
			group,ok:=decl.(*ast.GenDecl);if !ok||group.Tok!=token.TYPE{continue}
			for _,node:=range group.Specs {
				spec,ok:=node.(*ast.TypeSpec);if !ok{continue}
				iface,ok:=spec.Type.(*ast.InterfaceType);if !ok{continue}
				doc:=spec.Doc;if doc==nil{doc=group.Doc}
				receiver,present,err:=argument(doc,"Mapper");if err!=nil{return result,false,err}
				if !present{continue}
				if found{return result,false,fmt.Errorf("one @Mapper interface per service is supported")}
				if !token.IsIdentifier(receiver){return result,false,fmt.Errorf("invalid @Mapper receiver %q",receiver)}
				found=true
				result.receiver=receiver
				result.imports=file.Imports
				for _,member:=range iface.Methods.List {
					if len(member.Names)!=1{return result,false,fmt.Errorf("embedded mapper interfaces are not supported")}
					rules,err:=mappingRules(member.Doc);if err!=nil{return result,false,fmt.Errorf("%s: %w",member.Names[0].Name,err)}
					if len(rules)==0{continue}
					fn,ok:=member.Type.(*ast.FuncType)
					if !ok||fn.Params==nil||len(fn.Params.List)!=1||len(fn.Params.List[0].Names)!=1||fn.Results==nil||len(fn.Results.List)!=1 {
						return result,false,fmt.Errorf("%s: @Mapping requires one named input and one struct result",member.Names[0].Name)
					}
					destination,err:=destinationType(backend,file.Imports,fn.Results.List[0].Type)
					if err!=nil{return result,false,fmt.Errorf("%s: %w",member.Names[0].Name,err)}
					if len(rules)!=len(destination){return result,false,fmt.Errorf("%s: mapped %d fields, destination requires %d",member.Names[0].Name,len(rules),len(destination))}
					ordered:=make([]fieldMapping,0,len(destination))
					for _,field:=range destination {
						expr,ok:=rules[field]
						if !ok{return result,false,fmt.Errorf("%s: unmapped destination field %s",member.Names[0].Name,field)}
						ordered=append(ordered,fieldMapping{name:field,expression:expr})
					}
					var signature bytes.Buffer
					if err:=format.Node(&signature,set,fn);err!=nil{return result,false,err}
					var resultType bytes.Buffer
					if err:=format.Node(&resultType,set,fn.Results.List[0].Type);err!=nil{return result,false,err}
					result.methods=append(result.methods,mappingMethod{
						name:member.Names[0].Name,signature:strings.TrimPrefix(signature.String(),"func"),
						output:resultType.String(),fields:ordered,
					})
				}
			}
		}
	}
	return result,found,nil
}

func argument(doc *ast.CommentGroup,name string)(string,bool,error) {
	if doc==nil{return "",false,nil}
	var value string
	found:=false
	prefix:="@"+name
	for _,comment:=range doc.List {
		text:=strings.TrimSpace(strings.TrimPrefix(comment.Text,"//"))
		if text!=prefix&&!strings.HasPrefix(text,prefix+" "){continue}
		if found{return "",false,fmt.Errorf("duplicate %s",prefix)}
		found=true
		value=strings.TrimSpace(strings.TrimPrefix(text,prefix))
		if value==""||strings.ContainsAny(value," \t\r\n"){return "",false,fmt.Errorf("%s requires one identifier",prefix)}
	}
	return value,found,nil
}

func mappingRules(doc *ast.CommentGroup)(map[string]string,error) {
	rules:=map[string]string{}
	if doc==nil{return rules,nil}
	for _,comment:=range doc.List {
		line:=strings.TrimSpace(strings.TrimPrefix(comment.Text,"//"))
		if !strings.HasPrefix(line,"@Mapping "){continue}
		decl:=strings.TrimSpace(strings.TrimPrefix(line,"@Mapping "))
		sep:=strings.Index(decl,"=")
		if sep<=0{return nil,fmt.Errorf("expected @Mapping Destination=expression")}
		field,expression:=strings.TrimSpace(decl[:sep]),strings.TrimSpace(decl[sep+1:])
		if !token.IsIdentifier(field){return nil,fmt.Errorf("invalid destination field %q",field)}
		if _,exists:=rules[field];exists{return nil,fmt.Errorf("duplicate mapping for %s",field)}
		if expression==""{return nil,fmt.Errorf("missing expression for %s",field)}
		if _,err:=parser.ParseExpr(expression);err!=nil{return nil,fmt.Errorf("%s: invalid Go mapping expression: %w",field,err)}
		rules[field]=expression
	}
	return rules,nil
}

func destinationType(backend string,imports []*ast.ImportSpec,expr ast.Expr)([]string,error) {
	sel,ok:=expr.(*ast.SelectorExpr);if !ok{return nil,fmt.Errorf("mapped return type must be an imported DTO struct")}
	alias,ok:=sel.X.(*ast.Ident);if !ok{return nil,fmt.Errorf("unsupported destination type")}
	var importPath string
	for _,entry:=range imports {
		p,err:=strconv.Unquote(entry.Path.Value);if err!=nil{return nil,err}
		id:=path.Base(p);if entry.Name!=nil{id=entry.Name.Name}
		if id==alias.Name{importPath=p;break}
	}
	if !strings.HasPrefix(importPath,"grovehub/services/"){return nil,fmt.Errorf("destination %s is not a service DTO import",alias.Name)}
	dir:=filepath.Join(backend,strings.TrimPrefix(importPath,"grovehub/"))
	files,err:=filepath.Glob(filepath.Join(dir,"*.go"));if err!=nil{return nil,err}
	sort.Strings(files)
	for _,file:=range files {
		if strings.HasSuffix(file,"_test.go"){continue}
		astFile,err:=parser.ParseFile(token.NewFileSet(),file,nil,0);if err!=nil{return nil,err}
		for _,decl:=range astFile.Decls {
			group,ok:=decl.(*ast.GenDecl);if !ok||group.Tok!=token.TYPE{continue}
			for _,node:=range group.Specs {
				spec,ok:=node.(*ast.TypeSpec);if !ok||spec.Name.Name!=sel.Sel.Name{continue}
				typ,ok:=spec.Type.(*ast.StructType);if !ok{return nil,fmt.Errorf("destination %s is not a struct",sel.Sel.Name)}
				fields:=make([]string,0,len(typ.Fields.List))
				for _,field:=range typ.Fields.List {
					if len(field.Names)==0{return nil,fmt.Errorf("embedded fields are forbidden")}
					for _,name:=range field.Names{fields=append(fields,name.Name)}
				}
				return fields,nil
			}
		}
	}
	return nil,fmt.Errorf("destination struct %s not found in %s",sel.Sel.Name,dir)
}

func write(group mappingGroup,out string)error {
	var body strings.Builder
	for _,m:=range group.methods {
		fmt.Fprintf(&body,"func (mapper *%s) %s%s {\nreturn %s{\n",group.receiver,m.name,m.signature,m.output)
		for _,f:=range m.fields{fmt.Fprintf(&body,"%s: %s,\n",f.name,f.expression)}
		body.WriteString("}\n}\n\n")
	}
	var b strings.Builder
	b.WriteString("// Code generated by GrooveHub @Mapper. DO NOT EDIT.\npackage impl\n\n")
	// Only retain imports used in generated method signatures/bodies.
	b.WriteString("import (\n")
	for _,imp:=range group.imports {
		importPath,err:=strconv.Unquote(imp.Path.Value);if err!=nil{return err}
		name:=path.Base(importPath)
		if imp.Name!=nil{name=imp.Name.Name}
		if !strings.Contains(body.String(),name+"."){continue}
		if imp.Name!=nil{fmt.Fprintf(&b,"%s ",name)}
		fmt.Fprintf(&b,"%q\n",importPath)
	}
	b.WriteString(")\n\n")
	b.WriteString(body.String())
	formatted,err:=format.Source([]byte(b.String()));if err!=nil{return err}
	old,err:=os.ReadFile(out)
	if err==nil&&bytes.Equal(old,formatted){return nil}
	if err!=nil&&!os.IsNotExist(err){return err}
	if err:=os.MkdirAll(filepath.Dir(out),0755);err!=nil{return err}
	return os.WriteFile(out,formatted,0644)
}

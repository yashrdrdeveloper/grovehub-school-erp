// Command entities generates optional builders/getters/setters only when
// explicitly requested on NON-entity types. @Entity itself only maps a GORM
// table and must not imply any constructor or accessor generation.
// Code is written only to the owning service's ignored .generated/entity tree.
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

type entitySpec struct {
	name             string
	fields           []fieldSpec
	builder, getter, setter bool
}
type fieldSpec struct {
	name, typ string
}
type sourceFile struct {
	name, pkg string
	imports []*ast.ImportSpec
	entities []entitySpec
}
type methodIndex map[string]map[string]bool

func main() {
	if len(os.Args)!=2 {
		fmt.Fprintln(os.Stderr,"usage: entities <backend-dir>")
		os.Exit(2)
	}
	if err:=generateAll(os.Args[1]);err!=nil {
		fmt.Fprintln(os.Stderr,err)
		os.Exit(1)
	}
}

func generateAll(backend string) error {
	services,err:=os.ReadDir(filepath.Join(backend,"services"))
	if err!=nil{return err}
	for _,service:=range services {
		if !service.IsDir(){continue}
		owner:=filepath.Join(backend,"services",service.Name())
		dir:=filepath.Join(owner,"entity")
        if info,err:=os.Stat(filepath.Join(owner,"internal","app","entity"));err==nil&&info.IsDir(){
            dir=filepath.Join(owner,"internal","app","entity")
        }
		if info,err:=os.Stat(dir);err!=nil||!info.IsDir(){continue}
        rel, err := filepath.Rel(owner, dir)
        if err != nil { return err }
		if err:=generate(dir,filepath.Join(owner,".generated",rel));err!=nil {
			return fmt.Errorf("%s entities: %w",service.Name(),err)
		}
	}
	return nil
}

func generate(sourceDir,outputDir string) error {
	files,err:=filepath.Glob(filepath.Join(sourceDir,"*.go"))
	if err!=nil{return err}
	sort.Strings(files)
	type parsedFile struct {
		name string
		fset *token.FileSet
		tree *ast.File
	}
	var parsed []parsedFile
	methods:=methodIndex{}
	var packageName string
	for _,file:=range files {
		if strings.HasSuffix(file,"_test.go")||strings.HasSuffix(file,"_gombok.go")||strings.HasSuffix(file,".gen.go"){continue}
		fset:=token.NewFileSet()
		tree,err:=parser.ParseFile(fset,file,nil,parser.ParseComments)
		if err!=nil{return err}
		if packageName==""{packageName=tree.Name.Name}else if packageName!=tree.Name.Name{return fmt.Errorf("mixed packages in %s",sourceDir)}
		parsed=append(parsed,parsedFile{name:file,fset:fset,tree:tree})
		for _,decl:=range tree.Decls {
			fn,ok:=decl.(*ast.FuncDecl)
			if !ok||fn.Recv==nil||len(fn.Recv.List)!=1{continue}
			receiver:=fn.Recv.List[0].Type
			if ptr,ok:=receiver.(*ast.StarExpr);ok{receiver=ptr.X}
			id,ok:=receiver.(*ast.Ident);if !ok{continue}
			if methods[id.Name]==nil{methods[id.Name]=map[string]bool{}}
			methods[id.Name][fn.Name.Name]=true
		}
	}

	created:=map[string]bool{}
	for _,file:=range parsed {
		src:=sourceFile{name:filepath.Base(file.name),pkg:packageName,imports:file.tree.Imports}
		for _,decl:=range file.tree.Decls {
			group,ok:=decl.(*ast.GenDecl)
			if !ok||group.Tok!=token.TYPE{continue}
			for _,node:=range group.Specs {
				spec,ok:=node.(*ast.TypeSpec);if !ok{continue}
				typ,ok:=spec.Type.(*ast.StructType);if !ok{continue}
				doc:=spec.Doc;if doc==nil{doc=group.Doc}
				flags:=annotations(doc)
				// @Entity maps GORM tables only. Builder/getter/setter
				// generation is reserved for explicit non-entity directives.
				if flags["Entity"] {
					if flags["Builder"]||flags["Getter"]||flags["Setter"] {
						return fmt.Errorf("%s.%s: @Entity cannot request builders or accessors",src.name,spec.Name.Name)
					}
					continue
				}
				if !flags["Builder"]&&!flags["Getter"]&&!flags["Setter"]{continue}

				e:=entitySpec{
					name:spec.Name.Name,
					builder:flags["Builder"],
					getter:flags["Getter"],
					setter:flags["Setter"],
				}
				for _,field:=range typ.Fields.List {
					if len(field.Names)==0{return fmt.Errorf("%s.%s: embedded fields are not supported",src.name,e.name)}
					var rendered bytes.Buffer
					if err:=format.Node(&rendered,file.fset,field.Type);err!=nil{return err}
					for _,name:=range field.Names{
						if !name.IsExported(){return fmt.Errorf("%s.%s: unexported field %s cannot be generated",src.name,e.name,name.Name)}
						e.fields=append(e.fields,fieldSpec{name:name.Name,typ:rendered.String()})
					}
				}
				// Preserve the current Foundation Person semantics: custom value
				// getters/setters disable whole generated accessor categories.
				for method:=range methods[e.name] {
					if strings.HasPrefix(method,"Get"){e.getter=false}
					if strings.HasPrefix(method,"Set"){e.setter=false}
				}
				src.entities=append(src.entities,e)
			}
		}
		filename:=strings.TrimSuffix(src.name,".go")+"_gombok.go"
		output:=filepath.Join(outputDir,filename)
		if len(src.entities)==0 {
			if err:=os.Remove(output);err!=nil&&!os.IsNotExist(err){return err}
			continue
		}
		if err:=write(src,output);err!=nil{return err}
		created[filename]=true
	}
	// Remove files from annotations that were removed or renamed. Other
	// generated outputs (e.g. GORM table methods) are not touched.
	entries,err:=os.ReadDir(outputDir)
	if os.IsNotExist(err){return nil}
	if err!=nil{return err}
	for _,entry:=range entries{
		if entry.IsDir()||!strings.HasSuffix(entry.Name(),"_gombok.go")||created[entry.Name()]{continue}
		if err:=os.Remove(filepath.Join(outputDir,entry.Name()));err!=nil{return err}
	}
	return nil
}

func annotations(doc *ast.CommentGroup) map[string]bool {
	flags:=map[string]bool{}
	if doc==nil{return flags}
	for _,comment:=range doc.List {
		line:=strings.TrimSpace(strings.TrimPrefix(comment.Text,"//"))
		for _,name:=range []string{"Entity","Builder","Getter","Setter"} {
			key:="@"+name
			if line==key||strings.HasPrefix(line,key+" "){flags[name]=true}
		}
	}
	return flags
}

func write(src sourceFile,output string)error {
	var body strings.Builder
	for _,e:=range src.entities {
		if e.builder {
			fmt.Fprintf(&body,"type %sBuilder struct { target *%s }\n\n",e.name,e.name)
			for _,field:=range e.fields {
				fmt.Fprintf(&body,"func (b %sBuilder) With%s(value %s) %sBuilder { b.target.%s = value; return b }\n",
					e.name,field.name,field.typ,e.name,field.name)
			}
			fmt.Fprintf(&body,"func (b %sBuilder) Build() %s { return *b.target }\n",e.name,e.name)
			fmt.Fprintf(&body,"func New%sBuilder() %sBuilder { return %sBuilder{target: &%s{}} }\n\n",e.name,e.name,e.name,e.name)
		}
		for _,field:=range e.fields {
			if e.getter {
				fmt.Fprintf(&body,"func (value *%s) Get%s() %s { return value.%s }\n",e.name,field.name,field.typ,field.name)
			}
			if e.setter {
				fmt.Fprintf(&body,"func (value *%s) Set%s(next %s) { value.%s = next }\n",e.name,field.name,field.typ,field.name)
			}
		}
		body.WriteString("\n")
	}
	var b strings.Builder
	fmt.Fprintf(&b,"// Code generated by GrooveHub @Entity (Gombok-compatible). DO NOT EDIT.\npackage %s\n\n",src.pkg)
	var imports []string
	for _,entry:=range src.imports {
		importPath,err:=strconv.Unquote(entry.Path.Value);if err!=nil{return err}
		alias:=path.Base(importPath)
		if entry.Name!=nil{alias=entry.Name.Name}
		if !strings.Contains(body.String(),alias+"."){continue}
		prefix:=""
		if entry.Name!=nil{prefix=alias+" "}
		imports=append(imports,prefix+strconv.Quote(importPath))
	}
	if len(imports)>0 {
		b.WriteString("import (\n")
		for _,line:=range imports{b.WriteString(line+"\n")}
		b.WriteString(")\n\n")
	}
	b.WriteString(body.String())
	formatted,err:=format.Source([]byte(b.String()))
	if err!=nil{return fmt.Errorf("format %s: %w",src.name,err)}
	if old,err:=os.ReadFile(output);err==nil&&bytes.Equal(old,formatted){return nil}else if err!=nil&&!os.IsNotExist(err){return err}
	if err:=os.MkdirAll(filepath.Dir(output),0755);err!=nil{return err}
	return os.WriteFile(output,formatted,0644)
}


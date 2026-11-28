package helpers

import (
	"go/token"
	"grovehub/libs/core/annotation/internal/src/helpers/layout"
	"fmt"
	"go/ast"
	"go/parser"
	"os"
	"path/filepath"
	"strings"
)

func Generate(sourceDir,outputDir string) error {
	return GenerateWithRenderer(sourceDir,outputDir,layout.Render)
}

func GenerateWithRenderer(sourceDir,outputDir string,render layout.Renderer)error{
	files,err:=layout.GoFiles(sourceDir)
	if err!=nil{return err}
	type parsedFile struct {
		name string
		source layout.ParsedSource
		
	}
	var parsed []parsedFile
	methods:=methodIndex{}
	var packageName string
	for _,file:=range files {
		if strings.HasSuffix(file,"_test.go")||strings.HasSuffix(file,"_gombok.go")||strings.HasSuffix(file,".gen.go"){continue}
		source,err:=layout.ParseSource(file,parser.ParseComments)
		if err!=nil{return err}
		tree:=source.Tree
		if packageName==""{packageName=tree.Name.Name}else if packageName!=tree.Name.Name{return fmt.Errorf("mixed packages in %s",sourceDir)}
		parsed=append(parsed,parsedFile{name:file,source:source})
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
		src:=sourceFile{name:filepath.Base(file.name),pkg:packageName,imports:file.source.Imports}
		for _,decl:=range file.source.Tree.Decls {
			group,ok:=decl.(*ast.GenDecl)
			if !ok||group.Tok!=token.TYPE{continue}
			for _,spec:=range layout.TypeSpecs(group) {
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
					rendered,err:=render(file.source.Set,field.Type);if err!=nil{return err}
					for _,name:=range field.Names{
						if !name.IsExported(){return fmt.Errorf("%s.%s: unexported field %s cannot be generated",src.name,e.name,name.Name)}
						e.fields=append(e.fields,fieldSpec{name:name.Name,typ:rendered})
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

package helpers

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/helpers/layout"
)

func destinationType(backend,modulePath string,imports []*ast.ImportSpec,expr ast.Expr)([]string,error) {
	sel,ok:=expr.(*ast.SelectorExpr);if !ok{return nil,fmt.Errorf("mapped return type must be an imported DTO struct")}
	alias,ok:=sel.X.(*ast.Ident);if !ok{return nil,fmt.Errorf("unsupported destination type")}
	var importPath string
	for _,entry:=range imports {
		p,err:=strconv.Unquote(entry.Path.Value);if err!=nil{return nil,err}
		id:=path.Base(p);if entry.Name!=nil{id=entry.Name.Name}
		if id==alias.Name{importPath=p;break}
	}
	if !strings.HasPrefix(importPath,modulePath+"/services/"){return nil,fmt.Errorf("destination %s is not a service DTO import",alias.Name)}
	dir,err:=layout.ResolveImport(backend,modulePath,importPath)
	if err!=nil{return nil,err}
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

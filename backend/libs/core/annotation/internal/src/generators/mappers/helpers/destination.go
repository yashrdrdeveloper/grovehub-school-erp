package helpers

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"grovehub/libs/core/annotation/internal/src/helpers/layout"
)

func destinationType(backend,modulePath string,imports []layout.ImportBinding,expr ast.Expr)([]string,error) {
	sel,ok:=expr.(*ast.SelectorExpr);if !ok{return nil,fmt.Errorf("mapped return type must be an imported DTO struct")}
	alias,ok:=sel.X.(*ast.Ident);if !ok{return nil,fmt.Errorf("unsupported destination type")}
	var importPath string
	for _,entry:=range imports {
		if entry.Name==alias.Name{importPath=entry.Path;break}
	}
	if !strings.HasPrefix(importPath,modulePath+"/services/"){return nil,fmt.Errorf("destination %s is not a service DTO import",alias.Name)}
	dir,err:=layout.ResolveImport(backend,modulePath,importPath)
	if err!=nil{return nil,err}
	files,err:=layout.GoFiles(dir);if err!=nil{return nil,err}
	for _,file:=range files {
		if strings.HasSuffix(file,"_test.go"){continue}
		astFile,err:=parser.ParseFile(token.NewFileSet(),file,nil,0);if err!=nil{return nil,err}
		for _,decl:=range astFile.Decls {
			group,ok:=decl.(*ast.GenDecl);if !ok||group.Tok!=token.TYPE{continue}
			for _,spec:=range layout.TypeSpecs(group) {
				if spec.Name.Name!=sel.Sel.Name{continue}
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

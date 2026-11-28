package helpers

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

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

type mapEntry struct {
	name string
	value ast.Expr
}

func entries(node ast.Expr)([]mapEntry,error){
	lit,ok:=node.(*ast.CompositeLit);if !ok{return nil,fmt.Errorf("expected static map literal")}
	result:=make([]mapEntry,0,len(lit.Elts))
	for _,item:=range lit.Elts {
		pair,ok:=item.(*ast.KeyValueExpr);if !ok{return nil,fmt.Errorf("map requires keyed entries")}
		name,err:=stringExpr(pair.Key);if err!=nil{return nil,err}
		result=append(result,mapEntry{name:name,value:pair.Value})
	}
	return result,nil
}

func names(node ast.Expr)([]string,error){
	items,err:=entries(node)
	if err!=nil{return nil,err}
	result:=make([]string,0,len(items))
	for _,item:=range items{result=append(result,item.name)}
	sort.Strings(result)
	return result,nil
}

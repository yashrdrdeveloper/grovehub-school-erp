package layout

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"strconv"
)

// ImportBinding retains the original explicit alias when generating imports.
type ImportBinding struct {
	Name string
	Path string
	Explicit bool
}

// ParsedSource keeps the syntax tree and decoded literals together. Generators
// share the same validation and do not reparse imports or field tags later.
type ParsedSource struct {
	Tree *ast.File
	Set *token.FileSet
	Imports []ImportBinding
	Tags map[*ast.Field]string
}

func ParseSource(filename string, mode parser.Mode)(ParsedSource,error){
	set:=token.NewFileSet()
	tree,err:=parser.ParseFile(set,filename,nil,mode)
	if err!=nil{return ParsedSource{},err}
	return DecodeSource(tree,set)
}

func DecodeSource(tree *ast.File,set *token.FileSet)(ParsedSource,error){
	if tree==nil{return ParsedSource{},fmt.Errorf("missing source syntax tree")}
	result:=ParsedSource{Tree:tree,Set:set,Tags:make(map[*ast.Field]string)}
	for _,entry:=range tree.Imports{
		if entry==nil||entry.Path==nil{return ParsedSource{},fmt.Errorf("missing import path")}
		value,err:=strconv.Unquote(entry.Path.Value)
		if err!=nil{return ParsedSource{},err}
		binding:=ImportBinding{Name:path.Base(value),Path:value}
		if entry.Name!=nil{binding.Name=entry.Name.Name;binding.Explicit=true}
		result.Imports=append(result.Imports,binding)
	}
	var decodeErr error
	ast.Inspect(tree,func(node ast.Node)bool{
		if decodeErr!=nil{return false}
		field,ok:=node.(*ast.Field)
		if !ok||field.Tag==nil{return true}
		value,err:=strconv.Unquote(field.Tag.Value)
		if err!=nil{decodeErr=err;return false}
		result.Tags[field]=value
		return true
	})
	if decodeErr!=nil{return ParsedSource{},decodeErr}
	return result,nil
}

// TypeSpecs returns typed declarations, safely ignoring unrelated AST specs.
func TypeSpecs(group *ast.GenDecl)[]*ast.TypeSpec{
	var result []*ast.TypeSpec
	for _,node:=range group.Specs{
		if spec,ok:=node.(*ast.TypeSpec);ok{result=append(result,spec)}
	}
	return result
}

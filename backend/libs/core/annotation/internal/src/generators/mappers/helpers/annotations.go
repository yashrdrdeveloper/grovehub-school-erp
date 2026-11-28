package helpers

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

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

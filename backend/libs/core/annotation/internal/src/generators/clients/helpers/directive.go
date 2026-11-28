package helpers

import (
	"fmt"
	"go/ast"
	"strings"
)

func directive(doc *ast.CommentGroup,name string)(string,bool,error) {
	if doc==nil {return "",false,nil}
	key:="@"+name
	found:=false
	var value string
	for _,comment:=range doc.List {
		line:=strings.TrimSpace(strings.TrimPrefix(comment.Text,"//"))
		if line!=key&&!strings.HasPrefix(line,key+" ") {continue}
		if found {return "",true,fmt.Errorf("duplicate %s",key)}
		found=true
		value=strings.TrimSpace(strings.TrimPrefix(line,key))
		if value==""||strings.ContainsAny(value," \t\r\n"){return "",true,fmt.Errorf("%s expects one binding",key)}
	}
	return value,found,nil
}

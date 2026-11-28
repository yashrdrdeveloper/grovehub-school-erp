package helpers

import (
	"go/ast"
	"strings"
)

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

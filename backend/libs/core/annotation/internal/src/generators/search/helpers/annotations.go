package helpers

import (
	"fmt"
	"go/ast"
	"strings"
)

func hasDirective(doc *ast.CommentGroup, name string) bool {
	if doc == nil { return false }
	for _, comment := range doc.List {
		text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
		if text == "@"+name || strings.HasPrefix(text, "@"+name+" ") { return true }
	}
	return false
}

func directive(group *ast.CommentGroup, name string) string {
	if group==nil {return ""}
	for _, c:=range group.List {
		line:=strings.TrimSpace(strings.TrimPrefix(c.Text,"//"))
		if line=="@"+name {return ""}
		if strings.HasPrefix(line,"@"+name+" ") {return strings.TrimSpace(strings.TrimPrefix(line,"@"+name))}
	}
	return ""
}

func options(raw string, expected []string)(map[string]string,error) {
	r:=make(map[string]string,len(expected))
	for _,part:=range strings.Fields(raw) {
		p:=strings.SplitN(part,"=",2)
		if len(p)!=2||p[1]=="" {return nil,fmt.Errorf("invalid annotation argument %q",part)}
		valid:=false
		for _,key:=range expected {if p[0]==key {valid=true}}
		if !valid||r[p[0]]!="" {return nil,fmt.Errorf("unknown or repeated annotation option %s",p[0])}
		r[p[0]]=p[1]
	}
	for _,key:=range expected {if r[key]=="" {return nil,fmt.Errorf("missing %s option",key)}}
	return r,nil
}

func searchOptions(raw string) (map[string]string,error) {
	result:=map[string]string{}
	allowed:=map[string]bool{"impl":true,"entity":true,"config":true,"alias":true}
	for _,item:=range strings.Fields(raw) {
		pair:=strings.SplitN(item,"=",2)
		if len(pair)!=2||pair[1]==""||!allowed[pair[0]]||result[pair[0]]!="" {
			return nil,fmt.Errorf("invalid or repeated @Searchable option %q",item)
		}
		result[pair[0]]=pair[1]
	}
	if result["config"]==""||result["alias"]=="" {
		return nil,fmt.Errorf("@Searchable requires config=... alias=...")
	}
	return result,nil
}

func optionalFields(group *ast.CommentGroup,name string)([]string,error){
	raw:=directive(group,name)
	if raw=="" {return nil,nil}
	opts,err:=options(raw,[]string{"fields"})
	if err!=nil{return nil,fmt.Errorf("@%s: %w",name,err)}
	return strings.Split(opts["fields"],","),nil
}

func fieldsFrom(group *ast.CommentGroup, name string) ([]string,error) {
	args,err:=options(directive(group,name),[]string{"fields"})
	if err!=nil {return nil,fmt.Errorf("@%s: %w",name,err)}
	return strings.Split(args["fields"],","),nil
}

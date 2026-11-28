package helpers

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strconv"
	"strings"
)

var topicPattern=regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

var eventPattern=regexp.MustCompile(`^[a-z][a-z0-9_]*(\.[a-z][a-z0-9_]*)+\.v([1-9][0-9]*)$`)

func annotation(doc *ast.CommentGroup,name string)(eventDefinition,bool,error){
    def:=eventDefinition{name:name,typeConst:name+"Type"}
    if doc==nil{return def,false,nil}
    seen:=false
    for _,comment:=range doc.List {
        line:=strings.TrimSpace(strings.TrimPrefix(comment.Text,"//"))
        if !strings.HasPrefix(line,"@Event "){continue}
        if seen{return def,true,fmt.Errorf("%s has duplicate @Event declarations",name)}
        seen=true
        args:=strings.Fields(strings.TrimSpace(strings.TrimPrefix(line,"@Event ")))
        seenKeys:=map[string]bool{}
        for _,arg:=range args {
            key,value,ok:=strings.Cut(arg,"=")
            if !ok||value==""||seenKeys[key]{return def,true,fmt.Errorf("%s has invalid event argument %q",name,arg)}
            seenKeys[key]=true
            switch key {
            case "topic":def.topic=value
            case "type":def.eventType=value
            case "source":def.source=value
            case "typeConst":def.typeConst=value
            default:return def,true,fmt.Errorf("%s has unknown @Event key %q",name,key)
            }
        }
        if !topicPattern.MatchString(def.topic)||!eventPattern.MatchString(def.eventType){
            return def,true,fmt.Errorf("%s requires valid topic and .vN event type",name)
        }
        if def.source!=""&&!topicPattern.MatchString(def.source){return def,true,fmt.Errorf("%s has invalid source",name)}
        if !token.IsIdentifier(def.typeConst){return def,true,fmt.Errorf("%s has invalid type constant",name)}
        m:=eventPattern.FindStringSubmatch(def.eventType)
        version,err:=strconv.Atoi(m[2]);if err!=nil{return def,true,err}
        def.version=version
    }
    return def,seen,nil
}

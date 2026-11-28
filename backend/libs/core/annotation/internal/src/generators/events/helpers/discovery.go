package helpers

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"sort"
	"strings"
)

func Discover(dir string)([]eventDefinition,string,error){
    files,err:=filepath.Glob(filepath.Join(dir,"*.go"));if err!=nil{return nil,"",err}
    sort.Strings(files)
    var defs []eventDefinition
    var pkg string
    for _,file:=range files {
        if strings.HasSuffix(file,"_test.go"){continue}
        tree,err:=parser.ParseFile(token.NewFileSet(),file,nil,parser.ParseComments);if err!=nil{return nil,"",err}
        if pkg==""{pkg=tree.Name.Name}else if pkg!=tree.Name.Name{return nil,"",fmt.Errorf("mixed packages")}
        for _,decl:=range tree.Decls {
            group,ok:=decl.(*ast.GenDecl);if !ok||group.Tok!=token.TYPE{continue}
            for _,node:=range group.Specs {
                spec,ok:=node.(*ast.TypeSpec);if !ok{continue}
                doc:=spec.Doc;if doc==nil{doc=group.Doc}
                def,found,err:=annotation(doc,spec.Name.Name);if err!=nil{return nil,"",err}
                if found{defs=append(defs,def)}
            }
        }
    }
    // Only one Source constant may be emitted per package, but it may be
    // shared by multiple event types that declare the same source.
    source:=""
    for _,def:=range defs{
        if def.source==""{continue}
        if source!=""&&source!=def.source{return nil,"",fmt.Errorf("conflicting event sources: %s and %s",source,def.source)}
        source=def.source
    }
    sort.Slice(defs,func(i,j int)bool{return defs[i].name<defs[j].name})
    return defs,pkg,nil
}

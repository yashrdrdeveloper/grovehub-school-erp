// Command events generates stable, typed event identity constants from
// domain-owned @Event contracts. Publication, transactions, payloads and
// subscription semantics remain explicit service responsibilities.
package main

import (
    "bytes"
    "fmt"
    "go/ast"
    "go/format"
    "go/parser"
    "go/token"
    "io/fs"
    "os"
    "path/filepath"
    "regexp"
    "sort"
    "strconv"
    "strings"
)

type eventDefinition struct {
    name, topic, eventType, source, typeConst string
    version int
}

var topicPattern=regexp.MustCompile(`^[a-z][a-z0-9_]*$`)
var eventPattern=regexp.MustCompile(`^[a-z][a-z0-9_]*(\.[a-z][a-z0-9_]*)+\.v([1-9][0-9]*)$`)

func main(){
    if len(os.Args)!=2{fmt.Fprintln(os.Stderr,"usage: events <backend-dir>");os.Exit(2)}
    if err:=generateAll(os.Args[1]);err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(1)}
}

func generateAll(backend string)error{
    services,err:=os.ReadDir(filepath.Join(backend,"services"));if err!=nil{return err}
    for _,service:=range services {
        if !service.IsDir(){continue}
        owner:=filepath.Join(backend,"services",service.Name())
        root:=filepath.Join(owner,"events")
        if info,err:=os.Stat(filepath.Join(owner,"internal","app","events"));err==nil&&info.IsDir(){
            root=filepath.Join(owner,"internal","app","events")
        }
        if _,err:=os.Stat(root);os.IsNotExist(err){continue}else if err!=nil{return err}
        if err:=filepath.WalkDir(root,func(dir string,item fs.DirEntry,walkErr error)error{
            if walkErr!=nil{return walkErr}
            if !item.IsDir(){return nil}
            rel,err:=filepath.Rel(owner,dir);if err!=nil{return err}
            out:=filepath.Join(owner,".generated",rel,"zz_events.gen.go")
            defs,pkg,err:=discover(dir);if err!=nil{return fmt.Errorf("%s: %w",dir,err)}
            if len(defs)==0 {
                if err:=os.Remove(out);err!=nil&&!os.IsNotExist(err){return err}
                return nil
            }
            return write(out,pkg,defs)
        });err!=nil{return err}
    }
    return nil
}

func discover(dir string)([]eventDefinition,string,error){
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

func write(output,pkg string,defs []eventDefinition)error{
    var b strings.Builder
    b.WriteString("// Code generated by GrooveHub @Event. DO NOT EDIT.\npackage "+pkg+"\n\nconst (\n")
    lastSource:=""
    for _,def:=range defs{
        fmt.Fprintf(&b,"%sTopic = %q\n",def.name,def.topic)
        fmt.Fprintf(&b,"%s = %q\n",def.typeConst,def.eventType)
        fmt.Fprintf(&b,"%sVersion = %d\n",def.name,def.version)
        if def.source!=""&&def.source!=lastSource{
            fmt.Fprintf(&b,"Source = %q\n",def.source)
            lastSource=def.source
        }
    }
    b.WriteString(")\n")
    data,err:=format.Source([]byte(b.String()));if err!=nil{return err}
    previous,err:=os.ReadFile(output)
    if err==nil&&bytes.Equal(previous,data){return nil}
    if err!=nil&&!os.IsNotExist(err){return err}
    if err:=os.MkdirAll(filepath.Dir(output),0755);err!=nil{return err}
    return os.WriteFile(output,data,0644)
}

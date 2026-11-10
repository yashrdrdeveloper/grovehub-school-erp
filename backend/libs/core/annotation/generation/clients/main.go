// Command clients generates typed in-process forwarding adapters from
// @Client and @Proxy annotations on public client interfaces.
package main

import (
    "bytes"
    "fmt"
    "go/ast"
    "go/format"
    "go/parser"
    "go/token"
    "os"
    "path/filepath"
    "sort"
    "strings"
)

type proxyMethod struct{
    name, signature, delegate, arguments string
    hasResults bool
}
type clientProxy struct{
    receiver string
    imports []*ast.ImportSpec
    methods []proxyMethod
}

func main(){
    if len(os.Args)!=2{fmt.Fprintln(os.Stderr,"usage: clients <backend-dir>");os.Exit(2)}
    if err:=generateAll(os.Args[1]);err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(1)}
}

func generateAll(backend string)error {
    services,err:=os.ReadDir(filepath.Join(backend,"services"));if err!=nil{return err}
    for _,service:=range services{
        if !service.IsDir(){continue}
        owner:=filepath.Join(backend,"services",service.Name())
        clientDir:=filepath.Join(owner,"client")
        if info,err:=os.Stat(clientDir);err!=nil||!info.IsDir(){continue}
        proxies,err:=discover(clientDir);if err!=nil{return fmt.Errorf("%s: %w",service.Name(),err)}
        output:=filepath.Join(owner,".generated","client","impl","zz_clients.gen.go")
        if len(proxies)==0{
            if err:=os.Remove(output);err!=nil&&!os.IsNotExist(err){return err}
            continue
        }
        if err:=writeProxies(proxies,output);err!=nil{return fmt.Errorf("%s: %w",service.Name(),err)}
    }
    return nil
}

func discover(dir string)([]clientProxy,error){
    files,err:=filepath.Glob(filepath.Join(dir,"*.go"));if err!=nil{return nil,err}
    sort.Strings(files)
    var result []clientProxy
    for _,path:=range files{
        if strings.HasSuffix(path,"_test.go"){continue}
        set:=token.NewFileSet()
        file,err:=parser.ParseFile(set,path,nil,parser.ParseComments);if err!=nil{return nil,err}
        for _,decl:=range file.Decls{
            group,ok:=decl.(*ast.GenDecl);if !ok||group.Tok!=token.TYPE{continue}
            for _,node:=range group.Specs{
                spec,ok:=node.(*ast.TypeSpec);if !ok{continue}
                iface,ok:=spec.Type.(*ast.InterfaceType);if !ok{continue}
                doc:=spec.Doc;if doc==nil{doc=group.Doc}
                receiver,found,err:=directive(doc,"Client");if err!=nil{return nil,err}
                if !found{continue}
                if !token.IsIdentifier(receiver){return nil,fmt.Errorf("%s: invalid @Client receiver %q",spec.Name.Name,receiver)}
                proxy:=clientProxy{receiver:receiver,imports:file.Imports}
                for _,member:=range iface.Methods.List{
                    if len(member.Names)!=1{return nil,fmt.Errorf("annotated @Client %s cannot embed interfaces",spec.Name.Name)}
                    target,found,err:=directive(member.Doc,"Proxy");if err!=nil{return nil,err}
                    if !found{return nil,fmt.Errorf("%s.%s: missing @Proxy binding",spec.Name.Name,member.Names[0].Name)}
                    parts:=strings.Split(target,".")
                    if len(parts)!=2||!token.IsIdentifier(parts[0])||!token.IsIdentifier(parts[1]){return nil,fmt.Errorf("%s.%s: invalid @Proxy binding %q",spec.Name.Name,member.Names[0].Name,target)}
                    fn,ok:=member.Type.(*ast.FuncType);if !ok{return nil,fmt.Errorf("%s.%s: expected function type",spec.Name.Name,member.Names[0].Name)}
                    var args []string
                    for _,param:=range fn.Params.List{
                        if len(param.Names)==0{return nil,fmt.Errorf("%s.%s: all arguments must be named",spec.Name.Name,member.Names[0].Name)}
                        for _,name:=range param.Names{args=append(args,name.Name)}
                    }
                    var signature bytes.Buffer
                    if err:=format.Node(&signature,set,fn);err!=nil{return nil,err}
                    proxy.methods=append(proxy.methods,proxyMethod{
                        name:member.Names[0].Name,
                        signature:strings.TrimPrefix(signature.String(),"func"),
                        delegate:target,
                        arguments:strings.Join(args,", "),
                        hasResults:fn.Results!=nil&&len(fn.Results.List)>0,
                    })
                }
                result=append(result,proxy)
            }
        }
    }
    return result,nil
}

func directive(doc *ast.CommentGroup,name string)(string,bool,error){
    if doc==nil{return "",false,nil}
    key:="@"+name
    found:=false
    var value string
    for _,comment:=range doc.List{
        line:=strings.TrimSpace(strings.TrimPrefix(comment.Text,"//"))
        if line!=key&&!strings.HasPrefix(line,key+" "){continue}
        if found{return "",true,fmt.Errorf("duplicate %s",key)}
        found=true
        value=strings.TrimSpace(strings.TrimPrefix(line,key))
        if value==""||strings.ContainsAny(value," \t\r\n"){return "",true,fmt.Errorf("%s expects one binding",key)}
    }
    return value,found,nil
}

func writeProxies(proxies []clientProxy,output string)error{
    if len(proxies)!=1{return fmt.Errorf("one @Client interface per service is supported")}
    proxy:=proxies[0]
    var b strings.Builder
    b.WriteString("// Code generated by GrooveHub @Client/@Proxy. DO NOT EDIT.\npackage impl\n\nimport (\n")
    for _,imp:=range proxy.imports{
        if imp.Name!=nil{fmt.Fprintf(&b,"%s ",imp.Name.Name)}
        fmt.Fprintf(&b,"%s\n",imp.Path.Value)
    }
    b.WriteString(")\n\n")
    for _,m:=range proxy.methods{
        fmt.Fprintf(&b,"func (c *%s) %s%s {\n",proxy.receiver,m.name,m.signature)
        if m.hasResults{b.WriteString("return ")}
        fmt.Fprintf(&b,"c.%s(%s)\n}\n\n",m.delegate,m.arguments)
    }
    source,err:=format.Source([]byte(b.String()));if err!=nil{return err}
    previous,err:=os.ReadFile(output)
    if err==nil&&bytes.Equal(previous,source){return nil}
    if err!=nil&&!os.IsNotExist(err){return err}
    if err:=os.MkdirAll(filepath.Dir(output),0755);err!=nil{return err}
    return os.WriteFile(output,source,0644)
}

// Command inventory enumerates every maintained Foundation app function and
// reports missing mirrored unit-test files. It never invents test coverage.
package main

import (
    "encoding/json"
    "fmt"
    "go/ast"
    "go/parser"
    "go/token"
    "io/fs"
    "os"
    "path/filepath"
    "sort"
    "strings"
)

type Function struct {
    Source string `json:"source"`
    Name string `json:"name"`
    Mirror string `json:"mirror"`
    Mirrored bool `json:"mirrored"`
}

type Report struct {
    SourceFiles int `json:"sourceFiles"`
    TestFiles int `json:"testFiles"`
    Functions int `json:"functions"`
    MirroredFunctions int `json:"mirroredFunctions"`
    Missing []Function `json:"missing"`
}

func main() {
    if len(os.Args)<2||len(os.Args)>3 {fmt.Fprintln(os.Stderr,"usage: inventory <foundation-service-directory> [--strict]");os.Exit(2)}
    report,err:=inspect(os.Args[1])
    if err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(1)}
    output,err:=json.MarshalIndent(report,"","  ")
    if err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(1)}
    fmt.Println(string(output))
    if len(os.Args)==3&&os.Args[2]=="--strict"&&len(report.Missing)>0 {os.Exit(1)}
}

func inspect(owner string)(Report,error){
    root:=filepath.Join(owner,"internal","app")
    testRoot:=filepath.Join(owner,"internal","tests","unit-test")
    var result Report
    tracked:=map[string]bool{}
    err:=filepath.WalkDir(root,func(path string,entry fs.DirEntry,err error)error{
        if err!=nil{return err}
        if entry.IsDir(){return nil}
        if !strings.HasSuffix(path,".go")||strings.HasSuffix(path,"_test.go"){return nil}
        relative,err:=filepath.Rel(root,path)
        if err!=nil{return err}
        expected:=filepath.Join(testRoot,strings.TrimSuffix(relative,".go")+"_test.go")
        result.SourceFiles++
        if _,err:=os.Stat(expected);err==nil{tracked[expected]=true}
        f,err:=parser.ParseFile(token.NewFileSet(),path,nil,0)
        if err!=nil{return err}
        for _,decl:=range f.Decls{
            method,ok:=decl.(*ast.FuncDecl)
            if !ok{continue}
            name:=method.Name.Name
            if method.Recv!=nil&&len(method.Recv.List)>0{
                name="(receiver)."+name
            }
            result.Functions++
            if tracked[expected]{result.MirroredFunctions++}else{
                result.Missing=append(result.Missing,Function{Source:filepath.ToSlash(relative),Name:name,Mirror:filepath.ToSlash(expected),Mirrored:false})
            }
        }
        return nil
    })
    if err!=nil{return Report{},err}
    sort.Slice(result.Missing,func(i,j int)bool{
        if result.Missing[i].Source==result.Missing[j].Source{return result.Missing[i].Name<result.Missing[j].Name}
        return result.Missing[i].Source<result.Missing[j].Source
    })
    err=filepath.WalkDir(testRoot,func(path string,entry fs.DirEntry,err error)error{
        if os.IsNotExist(err){return nil}
        if err!=nil{return err}
        if !entry.IsDir()&&strings.HasSuffix(path,"_test.go"){result.TestFiles++}
        return nil
    })
    return result,err
}

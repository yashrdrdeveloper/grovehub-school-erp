// Package inventory enumerates every maintained service or core function and
// reports missing mirrored unit-test files. It never invents test coverage.
package inventory

import (

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



func Inspect(owner string)(Report,error){
    root:=filepath.Join(owner,"internal","app")
    if _,err:=os.Stat(root);os.IsNotExist(err){root=owner}
    testRoot:=filepath.Join(owner,"internal","tests","unit-test")
    if _,err:=os.Stat(testRoot);os.IsNotExist(err){testRoot=filepath.Join(owner,"tests","unit-test")}
    testRoots := []string{testRoot}
    if root == owner {
        entries, err := os.ReadDir(owner)
        if err != nil { return Report{}, err }
        for _, entry := range entries {
            if entry.IsDir() { testRoots = append(testRoots, filepath.Join(owner, entry.Name(), "internal", "resources", "tests", "unit-test")) }
        }
    }
    var result Report
    tracked:=map[string]bool{}
    err:=filepath.WalkDir(root,func(path string,entry fs.DirEntry,err error)error{
        if err!=nil{return err}
        if entry.IsDir(){
            if root==owner && (entry.Name()=="tests" || entry.Name()==".generated" || entry.Name()=="testdata") {
                return filepath.SkipDir
            }
            return nil
        }
        if !strings.HasSuffix(path,".go")||strings.HasSuffix(path,"_test.go"){return nil}
        relative,err:=filepath.Rel(root,path)
        if err!=nil{return err}
        expected:=filepath.Join(testRoot,strings.TrimSuffix(relative,".go")+"_test.go")
        parts := strings.Split(filepath.ToSlash(relative), "/")
        if root == owner && len(parts) > 3 && parts[1] == "internal" && parts[2] == "src" {
            expected = filepath.Join(owner, parts[0], "internal", "resources", "tests", "unit-test", filepath.FromSlash(strings.TrimSuffix(strings.Join(parts[3:], "/"), ".go")+"_test.go"))
        }
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
    for _, testRoot := range testRoots {
    err=filepath.WalkDir(testRoot,func(path string,entry fs.DirEntry,err error)error{
        if os.IsNotExist(err){return nil}
        if err!=nil{return err}
        if !entry.IsDir()&&strings.HasSuffix(path,"_test.go"){result.TestFiles++}
        return nil
    })
    if err != nil { return Report{}, err }
    }
    return result,nil
}

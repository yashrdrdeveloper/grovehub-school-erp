// Command overlay maps generated files in each module's .generated subtree
// into their maintained Go package locations without writing files there.
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/fs"
    "os"
    "path/filepath"
    "strings"
)

type configuration struct{
    Replace map[string]string `json:"Replace"`
}

func main() {
    if len(os.Args)!=2{fmt.Fprintln(os.Stderr,"usage: overlay <backend-directory>");os.Exit(2)}
    if err:=writeOverlay(os.Args[1]);err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(1)}
}

func writeOverlay(backend string)error {
    root,err:=filepath.Abs(backend);if err!=nil{return err}
    replacement:=map[string]string{}

    scan:=func(owner string)error{
        generatedRoot:=filepath.Join(owner,".generated")
        if _,err:=os.Stat(generatedRoot);os.IsNotExist(err){return nil}else if err!=nil{return err}
        return filepath.WalkDir(generatedRoot,func(path string,item fs.DirEntry,err error)error{
            if err!=nil{return err}
            if item.IsDir()||!strings.HasSuffix(item.Name(),".go"){return nil}
            if !item.Type().IsRegular(){return fmt.Errorf("generated Go file %s must be regular",path)}
            rel,err:=filepath.Rel(generatedRoot,path);if err!=nil{return err}
            logical:=filepath.Join(owner,rel)
            if _,err:=os.Stat(filepath.Dir(logical));err!=nil{return fmt.Errorf("generated source has no original package directory %s: %w",logical,err)}
            if _,err:=os.Stat(logical);err==nil{return fmt.Errorf("generated code conflicts with maintained source: %s",logical)}else if !os.IsNotExist(err){return err}
            if prior,ok:=replacement[logical];ok{return fmt.Errorf("duplicate virtual file %s: %s and %s",logical,prior,path)}
            replacement[logical]=path
            return nil
        })
    }

    services,err:=os.ReadDir(filepath.Join(root,"services"));if err!=nil{return err}
    for _,service:=range services {
        if !service.IsDir(){continue}
        if err:=scan(filepath.Join(root,"services",service.Name()));err!=nil{return err}
    }
    if err:=scan(filepath.Join(root,"libs","core"));err!=nil{return err}
    if len(replacement)==0{return fmt.Errorf("no generated Go files found; run generate-annotations.sh first")}

    outDir:=filepath.Join(root,".generated")
    if err:=os.MkdirAll(outDir,0755);err!=nil{return err}
    out:=filepath.Join(outDir,"overlay.json")
    data,err:=json.MarshalIndent(configuration{Replace:replacement},"","  ");if err!=nil{return err}
    data=append(data,'\n')
    old,err:=os.ReadFile(out);if err==nil&&bytes.Equal(old,data){return nil}
    if err!=nil&&!os.IsNotExist(err){return err}
    temp,err:=os.CreateTemp(outDir,"overlay-*.json");if err!=nil{return err}
    defer os.Remove(temp.Name())
    if _,err=temp.Write(data);err!=nil{temp.Close();return err}
    if err=temp.Close();err!=nil{return err}
    return os.Rename(temp.Name(),out)
}

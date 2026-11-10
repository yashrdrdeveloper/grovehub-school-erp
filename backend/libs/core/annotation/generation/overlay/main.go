// Command overlay exposes service-local generated code as files inside the
// original Go packages without writing generated sources into maintained dirs.
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

type configuration struct {
    Replace map[string]string `json:"Replace"`
}

func main() {
    if len(os.Args)!=2 {fmt.Fprintln(os.Stderr,"usage: overlay <backend-directory>");os.Exit(2)}
    if err:=writeOverlay(os.Args[1]);err!=nil{fmt.Fprintln(os.Stderr,err);os.Exit(1)}
}

func writeOverlay(backend string)error {
    root,err:=filepath.Abs(backend);if err!=nil{return err}
    serviceFiles,err:=filepath.Glob(filepath.Join(root,"services","*",".generated","*","*.go"))
    if err!=nil{return err}
    coreFiles,err:=filepath.Glob(filepath.Join(root,"libs","core",".generated","*","*.go"))
    if err!=nil{return err}
    replacements:=make(map[string]string)
    for _,generated:=range append(serviceFiles,coreFiles...) {
        packageDir:=filepath.Base(filepath.Dir(generated))
        ownerDir:=filepath.Dir(filepath.Dir(filepath.Dir(generated)))
        sourceDir:=filepath.Join(ownerDir,packageDir)
        if !strings.HasPrefix(sourceDir,root+string(filepath.Separator)){return fmt.Errorf("invalid source location %s",sourceDir)}
        if info,err:=os.Stat(sourceDir);err!=nil||!info.IsDir(){return fmt.Errorf("generated %s has no maintained package at %s",generated,sourceDir)}
        logical:=filepath.Join(sourceDir,filepath.Base(generated))
        if _,err:=os.Stat(logical);err==nil{return fmt.Errorf("generated %s conflicts with maintained file %s",generated,logical)}else if !os.IsNotExist(err){return err}
        if previous,ok:=replacements[logical];ok{return fmt.Errorf("duplicate virtual file %s from %s and %s",logical,previous,generated)}
        replacements[logical]=generated
    }
    if len(replacements)==0 {return fmt.Errorf("no generated files found; run generate-entities.sh first")}
    outputDir:=filepath.Join(root,".generated")
    if err:=os.MkdirAll(outputDir,0755);err!=nil{return err}
    output:=filepath.Join(outputDir,"overlay.json")
    data,err:=json.MarshalIndent(configuration{Replace:replacements},"","  ");if err!=nil{return err}
    data=append(data,'\n')
    before,err:=os.ReadFile(output)
    if err==nil&&bytes.Equal(before,data){return nil}
    if err!=nil&&!os.IsNotExist(err){return err}
    temp,err:=os.CreateTemp(outputDir,"overlay-*.json");if err!=nil{return err}
    defer os.Remove(temp.Name())
    if _,err=temp.Write(data);err!=nil{temp.Close();return err}
    if err=temp.Close();err!=nil{return err}
    return os.Rename(temp.Name(),output)
}

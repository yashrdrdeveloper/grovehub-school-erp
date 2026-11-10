// Command overlay creates a Go build overlay from service-local generated sources.
package main

import (
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
 if len(os.Args)!=2 {
  fmt.Fprintln(os.Stderr,"usage: overlay <backend-directory>")
  os.Exit(2)
 }
 if err:=writeOverlay(os.Args[1]);err!=nil {fmt.Fprintln(os.Stderr,err);os.Exit(1)}
}

func writeOverlay(backend string) error {
 backend,err:=filepath.Abs(backend);if err!=nil{return err}
 config:=configuration{Replace:make(map[string]string)}
 serviceDirs,err:=filepath.Glob(filepath.Join(backend,"services","*",".generated","entity","*.go"))
 if err!=nil{return err}
 coreDirs,err:=filepath.Glob(filepath.Join(backend,"libs","core",".generated","entity","*.go"))
 if err!=nil{return err}
 for _,generated:=range append(serviceDirs,coreDirs...) {
  var sourceDir string
  if strings.Contains(generated,string(filepath.Separator)+"services"+string(filepath.Separator)) {
   sourceDir=filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(generated))),"entity")
  } else {
   sourceDir=filepath.Join(backend,"libs","core","entity")
  }
  if _,err:=os.Stat(sourceDir);err!=nil{return fmt.Errorf("missing source package for %s: %w",generated,err)}
  logical:=filepath.Join(sourceDir,filepath.Base(generated))
  if _,err:=os.Stat(logical);err==nil {
   return fmt.Errorf("generated file conflicts with checked-in source: %s; delete the old generated source file",logical)
  } else if !os.IsNotExist(err) {return err}
  config.Replace[logical]=generated
 }
 if len(config.Replace)==0 {return fmt.Errorf("no generated entity files found; run generate-entities.sh first")}
 outputDir:=filepath.Join(backend,".generated")
 if err:=os.MkdirAll(outputDir,0755);err!=nil{return err}
 output:=filepath.Join(outputDir,"overlay.json")
 data,err:=json.MarshalIndent(config,"","  ");if err!=nil{return err}
 data=append(data,'\n')
 temp,err:=os.CreateTemp(outputDir,"overlay-*.json")
 if err!=nil{return err}
 defer os.Remove(temp.Name())
 if _,err=temp.Write(data);err!=nil {temp.Close();return err}
 if err=temp.Close();err!=nil{return err}
 return os.Rename(temp.Name(),output)
}

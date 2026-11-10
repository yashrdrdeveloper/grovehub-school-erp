package main

import (
    "os"
    "path/filepath"
    "strings"
    "testing"
)

func TestRouteGeneration(t *testing.T) {
    backend:=t.TempDir()
    controllerDir:=filepath.Join(backend,"services","foundation","controller")
    if err:=os.MkdirAll(controllerDir,0755);err!=nil{t.Fatal(err)}
    src:=`package controller
import "net/http"
// @Routes /foundation/persons
type PersonController interface {
    // @POST /
    Create(http.ResponseWriter, *http.Request)
    // @GET /{gid}
    GetByGID(http.ResponseWriter, *http.Request)
}
`
    if err:=os.WriteFile(filepath.Join(controllerDir,"person.go"),[]byte(src),0644);err!=nil{t.Fatal(err)}
    if err:=generateAll(backend);err!=nil{t.Fatal(err)}
    output:=filepath.Join(backend,"services","foundation",".generated","routes","zz_routes.gen.go")
    first,err:=os.ReadFile(output);if err!=nil{t.Fatal(err)}
    for _,part:=range []string{"POST /foundation/persons", "GET /foundation/persons/{gid}", "registerPersonRoutes", "ServiceHandler"} {
        if !strings.Contains(string(first),part){t.Errorf("missing %s in generated route",part)}
    }
    if err:=generateAll(backend);err!=nil{t.Fatal(err)}
    again,err:=os.ReadFile(output);if err!=nil{t.Fatal(err)}
    if string(first)!=string(again){t.Fatal("generated routes are not deterministic")}
}

func TestRejectDuplicateRoute(t *testing.T) {
    dir:=t.TempDir()
    source:=`package controller
import "net/http"
// @Routes /persons
type PersonController interface {
    // @GET /
    A(http.ResponseWriter, *http.Request)
    // @GET /
    B(http.ResponseWriter, *http.Request)
}
`
    if err:=os.WriteFile(filepath.Join(dir,"a.go"),[]byte(source),0644);err!=nil{t.Fatal(err)}
    _,err:=discover(dir)
    if err==nil||!strings.Contains(err.Error(),"duplicate route"){t.Fatalf("expected duplicate route error, got %v",err)}
}

func TestRejectInvalidSignatureAndMissingVerb(t *testing.T) {
    for _,method:=range []string{"    // @GET /\n    A(string)", "    A(http.ResponseWriter, *http.Request)"} {
        dir:=t.TempDir()
        source:="package controller\nimport \"net/http\"\n// @Routes /persons\ntype PersonController interface {\n"+method+"\n}\n"
        if err:=os.WriteFile(filepath.Join(dir,"bad.go"),[]byte(source),0644);err!=nil{t.Fatal(err)}
        if _,err:=discover(dir);err==nil{t.Errorf("accepted invalid route declaration %s",method)}
    }
}

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


func TestCurrentBackendRoutesMatchContract(t *testing.T) {
    expected:=map[string][]string{
        "foundation":{
            "POST /foundation/persons",
            "POST /foundation/persons/search",
            "GET /foundation/persons/{gid}",
            "PATCH /foundation/persons/{gid}",
            "PATCH /foundation/persons/{gid}/deactivate",
        },
        "identity":{"POST /auth/register"},
        "integration":{
            "POST /integrations/connections",
            "GET /integrations/connections",
            "GET /integrations/connections/{id}",
            "DELETE /integrations/connections/{id}",
            "POST /integrations/connections/{id}/oauth/start",
            "GET /integrations/connections/oauth/{provider}/callback",
            "POST /integrations/connections/{id}/test",
            "GET /integrations/connections/{id}/tests",
            "POST /integrations/connections/{id}/watch",
            "POST /integrations/connections/{id}/actions/{action}",
            "POST /integrations/webhooks/gmail",
            "GET /integrations/providers",
            "GET /integrations/providers/{provider}",
        },
    }
    for service,patterns:=range expected{
        source:=filepath.Join("..","..","..","..","..","services",service,"controller")
        groups,err:=discover(source)
        if err!=nil{t.Fatalf("%s: %v",service,err)}
        actual:=map[string]bool{}
        for _,group:=range groups{
            for _,route:=range group.routes{
                actual[route.verb+" "+route.path]=true
            }
        }
        if len(actual)!=len(patterns){t.Errorf("%s expected %d routes, found %d: %v",service,len(patterns),len(actual),actual)}
        for _,p:=range patterns{
            if !actual[p]{t.Errorf("%s missing route %s",service,p)}
        }
    }
}

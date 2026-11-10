package main

import (
    "os"
    "path/filepath"
    "strings"
    "testing"
)

func TestAcceptedRoutesAndEntity(t *testing.T) {
    dir:=t.TempDir()
    path:=filepath.Join(dir,"source.go")
    src:=`package service
import "net/http"
// @Routes /people
type PersonController interface {
    // @GET /{id}
    Get(http.ResponseWriter,*http.Request)
}
// @Entity foundation.person_contacts
type Contact struct{ID string}
`
    if err:=os.WriteFile(path,[]byte(src),0644);err!=nil{t.Fatal(err)}
    if err:=CheckSource(path);err!=nil{t.Fatal(err)}
}

func TestRejectUnimplementedSecurityAnnotation(t *testing.T) {
    dir:=t.TempDir()
    path:=filepath.Join(dir,"source.go")
    src:=`package service
// @Authorize school:write
func Update() {}
`
    if err:=os.WriteFile(path,[]byte(src),0644);err!=nil{t.Fatal(err)}
    if err:=CheckSource(path);err==nil||!strings.Contains(err.Error(),"metadata only"){t.Fatalf("unsafe @Authorize accepted: %v",err)}
}

func TestRejectInvalidTarget(t *testing.T){
    dir:=t.TempDir()
    path:=filepath.Join(dir,"source.go")
    src:=`package service
// @Mapper Mapper
type Thing struct{}
`
    if err:=os.WriteFile(path,[]byte(src),0644);err!=nil{t.Fatal(err)}
    if err:=CheckSource(path);err==nil||!strings.Contains(err.Error(),"cannot annotate target"){t.Fatalf("incorrect mapper target accepted: %v",err)}
}

func TestCurrentBackendHasNoInertAnnotations(t *testing.T) {
    backend:=filepath.Join("..","..","..","..","..")
    if err:=checkServices(backend);err!=nil{t.Fatal(err)}
}

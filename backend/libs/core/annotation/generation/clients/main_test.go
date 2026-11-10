package main

import (
    "os"
    "path/filepath"
    "strings"
    "testing"
)

func TestClientProxyGeneration(t *testing.T){
    root:=t.TempDir()
    dir:=filepath.Join(root,"services","foundation","client")
    if err:=os.MkdirAll(dir,0755);err!=nil{t.Fatal(err)}
    source:=`package client
import "context"
// @Client personClient
type Client interface {
 // @Proxy personService.Get
 GetPerson(ctx context.Context, id string) (string,error)
}
`
    if err:=os.WriteFile(filepath.Join(dir,"person.go"),[]byte(source),0644);err!=nil{t.Fatal(err)}
    if err:=generateAll(root);err!=nil{t.Fatal(err)}
    output:=filepath.Join(root,"services","foundation",".generated","client","impl","zz_clients.gen.go")
    first,err:=os.ReadFile(output);if err!=nil{t.Fatal(err)}
    if !strings.Contains(string(first),"return c.personService.Get(ctx, id)") {t.Fatalf("bad proxy generation: %s",first)}
    if !strings.Contains(string(first),"GetPerson(ctx context.Context, id string) (string, error)") {t.Fatalf("wrong signature: %s",first)}
    if err:=generateAll(root);err!=nil{t.Fatal(err)}
    again,err:=os.ReadFile(output);if err!=nil{t.Fatal(err)}
    if string(first)!=string(again){t.Fatal("proxy generation is not deterministic")}
}

func TestCurrentServiceClientContracts(t *testing.T){
    expected:=map[string]int{"foundation":5,"integration":7}
    for service,count:=range expected {
        source:=filepath.Join("..","..","..","..","..","services",service,"client")
        proxies,err:=discover(source)
        if err!=nil{t.Fatalf("%s: %v",service,err)}
        if len(proxies)!=1{t.Fatalf("%s: expected one annotated interface, found %d",service,len(proxies))}
        if len(proxies[0].methods)!=count{t.Fatalf("%s: expected %d typed methods, found %d",service,count,len(proxies[0].methods))}
    }
}

func TestRejectMissingProxyBinding(t *testing.T){
    dir:=t.TempDir()
    source:=`package client
type Client interface {
 // @Client is not on this method
 Get() error
}
`
    source=strings.Replace(source,"type Client interface","// @Client clientImpl\ntype Client interface",1)
    if err:=os.WriteFile(filepath.Join(dir,"client.go"),[]byte(source),0644);err!=nil{t.Fatal(err)}
    _,err:=discover(dir)
    if err==nil||!strings.Contains(err.Error(),"missing @Proxy"){t.Fatalf("expected missing proxy error, got %v",err)}
}

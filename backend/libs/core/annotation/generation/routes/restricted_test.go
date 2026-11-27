package main

import (
 "os"
 "path/filepath"
 "testing"
)

func TestRestrictedProjectRouteURLs(t *testing.T) {
 root:=t.TempDir()
 dir:=filepath.Join(root,"person","protected")
 if err:=os.MkdirAll(dir,0755);err!=nil {t.Fatal(err)}
 source:=`package protected
import "net/http"
// @Routes /persons
type PersonController interface {
 // @GET /{gid}
 Get(http.ResponseWriter, *http.Request)
}
`
 if err:=os.WriteFile(filepath.Join(dir,"person.go"),[]byte(source),0644);err!=nil {t.Fatal(err)}
 groups,err:=discoverWithTier(root,"/foundation",true,true)
 if err!=nil {t.Fatal(err)}
 if len(groups)!=1 || len(groups[0].routes)!=1 || groups[0].routes[0].path!="/foundation/restricted/persons/{gid}" || groups[0].tier!="protected" {t.Fatalf("unexpected restricted groups: %+v",groups)}
 normal,err:=discover(root,"/foundation",true)
 if err!=nil {t.Fatal(err)}
 if normal[0].routes[0].path!="/foundation/protected/persons/{gid}" {t.Fatal("default generation changed")}
 other:=filepath.Join(root,"person","public")
 if err:=os.MkdirAll(other,0755);err!=nil {t.Fatal(err)}
 source="package public"+source[len("package protected"):]
 if err:=os.WriteFile(filepath.Join(other,"person.go"),[]byte(source),0644);err!=nil {t.Fatal(err)}
 if _,err:=discoverWithTier(root,"/foundation",true,true);err==nil {t.Fatal("duplicate restricted route accepted")}
}

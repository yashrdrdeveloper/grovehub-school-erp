package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMapperGenerationAndStrictFields(t *testing.T) {
	backend:=t.TempDir()
	dir:=filepath.Join(backend,"services","identity","mappers")
	outDir:=filepath.Join(backend,"services","identity","dto","response")
	if err:=os.MkdirAll(dir,0755);err!=nil{t.Fatal(err)}
	if err:=os.MkdirAll(outDir,0755);err!=nil{t.Fatal(err)}
	input:=`package mappers
import (
  response "grovehub/services/identity/dto/response"
  "grovehub/services/identity/entity"
)
// @Mapper AuthMapper
type AuthMapper interface {
  // @Mapping Email=account.GetEmail()
  // @Mapping ID=string(account.GetAccountId())
  ToResponse(account entity.Account) response.Account
}
`
	dto:=`package response
type Account struct {
 ID string
 Email string
}
`
	if err:=os.WriteFile(filepath.Join(dir,"auth.go"),[]byte(input),0644);err!=nil{t.Fatal(err)}
	if err:=os.WriteFile(filepath.Join(outDir,"account.go"),[]byte(dto),0644);err!=nil{t.Fatal(err)}
	if err:=generateAll(backend);err!=nil{t.Fatal(err)}
	out:=filepath.Join(backend,"services","identity",".generated","mappers","impl","zz_mappers.gen.go")
	got,err:=os.ReadFile(out);if err!=nil{t.Fatal(err)}
	for _,fragment:=range []string{"func (mapper *AuthMapper) ToResponse", "ID: string(account.GetAccountId())", "Email: account.GetEmail()"} {
		if !strings.Contains(string(got),fragment){t.Errorf("missing generated fragment %q:\n%s",fragment,got)}
	}
	before:=string(got)
	if err:=generateAll(backend);err!=nil{t.Fatal(err)}
	after,err:=os.ReadFile(out);if err!=nil{t.Fatal(err)}
	if string(after)!=before{t.Error("generation not deterministic")}
}

func TestRejectUnmappedFields(t *testing.T) {
	backend:=t.TempDir()
	dir:=filepath.Join(backend,"services","identity","mappers")
	outDir:=filepath.Join(backend,"services","identity","dto","response")
	if err:=os.MkdirAll(dir,0755);err!=nil{t.Fatal(err)}
	if err:=os.MkdirAll(outDir,0755);err!=nil{t.Fatal(err)}
	if err:=os.WriteFile(filepath.Join(outDir,"account.go"),[]byte("package response\ntype Account struct { ID string; Email string }\n"),0644);err!=nil{t.Fatal(err)}
	source:=`package mappers
import response "grovehub/services/identity/dto/response"
// @Mapper AuthMapper
type AuthMapper interface {
 // @Mapping ID=account.ID
 ToResponse(account struct{ID string}) response.Account
}
`
	if err:=os.WriteFile(filepath.Join(dir,"auth.go"),[]byte(source),0644);err!=nil{t.Fatal(err)}
	_,_,err:=discover(backend,dir)
	if err==nil||!strings.Contains(err.Error(),"mapped 1 fields"){t.Fatalf("expected unmapped-field error, got %v",err)}
}

func TestCurrentMappersIncludeExplicitRules(t *testing.T) {
	backend:=filepath.Join("..","..","..","..","..","..")
	for _,service:=range []string{"foundation","identity"} {
		dir:=filepath.Join(backend,"services",service,"mappers")
		m,found,err:=discover(backend,dir)
		if err!=nil{t.Fatalf("%s: %v",service,err)}
		if !found||len(m.methods)==0{t.Errorf("%s has no generated mapper methods",service)}
	}
}

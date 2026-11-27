package platform

import (
 "crypto/ed25519"
 "crypto/rand"
 "encoding/base64"
 "encoding/json"
 "testing"

 "github.com/google/uuid"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/services/workflow"
)

func TestSchoolGrantConfiguration(t *testing.T){
 makeGrant:=func()SchoolGrant{
  _,private,err:=ed25519.GenerateKey(rand.Reader);if err!=nil{t.Fatal(err)}
  return SchoolGrant{EnvironmentCode:"dev",TemplateID:uuid.NewString(),SchoolID:uuid.NewString(),Operation:"audit.activity-check",Origin:"http://127.0.0.1:8080",KeyID:uuid.NewString(),PrivateKey:base64.RawURLEncoding.EncodeToString(private),Requesters:[]string{uuid.NewString()},Reviewers:[]string{uuid.NewString()},Executors:[]string{uuid.NewString()}}
 }
 a,b:=makeGrant(),makeGrant()
 parse:=func(grants ...SchoolGrant)(map[workflow.SchoolBindingKey]workflow.SchoolBinding,error){
  raw,err:=json.Marshal(grants);if err!=nil{t.Fatal(err)}
  bindings,origins,err:=ParseSchoolGrants(string(raw),SchoolConfig{},nil)
  if err==nil&&len(origins)!=1{t.Fatal("origin dedup failed")}
  return bindings,err
 }
 bindings,err:=parse(a,b);if err!=nil||len(bindings)!=2{t.Fatal("multi-school grants failed",err)}
 svc:=workflow.NewWithSchoolBindings(nil,nil,bindings)
 views:=svc.SchoolGrants(entity.Operator{ID:a.Requesters[0],Persona:entity.L2})
 if len(views)!=1||views[0].SchoolID!=a.SchoolID||!views[0].CanRequest||views[0].CanReview||views[0].CanExecute{t.Fatal("role isolation failed",views)}
 if len(svc.SchoolGrants(entity.Operator{ID:uuid.NewString(),Persona:entity.L2}))!=0{t.Fatal("unassigned operator sees grant")}
 if len(svc.SchoolGrants(entity.Operator{ID:a.Reviewers[0],Persona:entity.L1}))!=0{t.Fatal("L1 sees review-only grant")}
 for _,bad:=range []SchoolGrant{
  a,
  func()SchoolGrant{v:=b;v.KeyID=a.KeyID;return v}(),
  func()SchoolGrant{v:=b;v.PrivateKey=a.PrivateKey;return v}(),
  func()SchoolGrant{v:=b;v.Operation="identity.mutate";return v}(),
  func()SchoolGrant{v:=b;v.Requesters=nil;return v}(),
  func()SchoolGrant{v:=b;v.Origin="http://127.0.0.1:9999";return v}(),
 }{if _,err:=parse(a,bad);err==nil{t.Fatal("invalid grant accepted")}}
 raw,_:=json.Marshal([]SchoolGrant{a})
 if _,_,err:=ParseSchoolGrants(string(raw),SchoolConfig{SchoolID:a.SchoolID},nil);err==nil{t.Fatal("legacy and new configuration mixed")}
 if _,_,err:=ParseSchoolGrants(string(raw),SchoolConfig{},[]ReservedCredential{{PrivateKey:a.PrivateKey}});err==nil{t.Fatal("reserved key reused")}
}

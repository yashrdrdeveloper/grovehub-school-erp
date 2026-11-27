package platform

import (
 "encoding/base64"
 "encoding/json"
 "errors"
 "io"
 "sort"
 "strings"

 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/services/workflow"
)

// SchoolGrant is deployment-owned configuration, not a registry entry that an
// operator can modify. Separate keys are required for every binding.
type SchoolGrant struct {
 EnvironmentCode string `json:"environmentCode"`
 TemplateID string `json:"templateId"`
 SchoolID string `json:"schoolId"`
 Operation string `json:"operation"`
 Origin string `json:"origin"`
 KeyID string `json:"keyId"`
 PrivateKey string `json:"privateKey"`
 Requesters []string `json:"requesters"`
 Reviewers []string `json:"reviewers"`
 Executors []string `json:"executors"`
}
type ReservedCredential struct{KeyID,PrivateKey string}

func ParseSchoolGrants(raw string,legacy SchoolConfig,reserved []ReservedCredential)(map[workflow.SchoolBindingKey]workflow.SchoolBinding,[]string,error){
 var grants []SchoolGrant
 if strings.TrimSpace(raw)!="" {
  if legacy!=(SchoolConfig{})||len(raw)>256*1024{return nil,nil,errors.New("invalid or conflicting school grant configuration")}
  d:=json.NewDecoder(strings.NewReader(raw));d.DisallowUnknownFields()
  var extra any
  if d.Decode(&grants)!=nil||grants==nil||d.Decode(&extra)!=io.EOF||len(grants)>100{return nil,nil,errors.New("invalid school grants")}
 }else if legacy!=(SchoolConfig{}){
  ids:=strings.Split(legacy.OperatorIDs,",")
  grants=[]SchoolGrant{{EnvironmentCode:legacy.EnvironmentCode,TemplateID:legacy.TemplateID,SchoolID:legacy.SchoolID,Operation:"audit.activity-check",Origin:legacy.Origin,KeyID:legacy.KeyID,PrivateKey:legacy.PrivateKey,Requesters:ids,Reviewers:ids,Executors:ids}}
 }
 bindings:=map[workflow.SchoolBindingKey]workflow.SchoolBinding{}
 origins:=map[string]bool{}
 usedIDs:=map[string]bool{}
 usedKeys:=map[string]bool{}
 for _,r:=range reserved{
  if r.KeyID!=""{usedIDs[r.KeyID]=true}
  if b,err:=base64.RawURLEncoding.DecodeString(r.PrivateKey);err==nil&&len(b)==64{usedKeys[string(b[32:])]=true}
 }
 for _,g:=range grants{
  if g.Operation!="audit.activity-check"{return nil,nil,errors.New("unsupported school operation")}
  key,err:=base64.RawURLEncoding.DecodeString(g.PrivateKey)
  if err!=nil||len(key)!=64||usedIDs[g.KeyID]||usedKeys[string(key[32:])]{return nil,nil,errors.New("school grant credential must be independent")}
  usedIDs[g.KeyID]=true;usedKeys[string(key[32:])]=true
  union:=map[string]bool{}
  for _,role:=range [][]string{g.Requesters,g.Reviewers,g.Executors}{
   if len(role)==0||len(role)>100{return nil,nil,errors.New("invalid school operator role")}
   seen:=map[string]bool{}
   for _,id:=range role{if !canonicalID(id)||seen[id]{return nil,nil,errors.New("invalid school operator role")};seen[id]=true;union[id]=true}
  }
  ids:=make([]string,0,len(union));for id:=range union{ids=append(ids,id)};sort.Strings(ids)
  one,err:=SchoolBindings(SchoolConfig{EnvironmentCode:g.EnvironmentCode,TemplateID:g.TemplateID,SchoolID:g.SchoolID,Origin:g.Origin,KeyID:g.KeyID,PrivateKey:g.PrivateKey,OperatorIDs:strings.Join(ids,",")})
  if err!=nil{return nil,nil,err}
  for k,b:=range one{
   if _,ok:=bindings[k];ok{return nil,nil,errors.New("duplicate school execution binding")}
   b.Requesters=g.Requesters;b.Reviewers=g.Reviewers;b.Executors=g.Executors
   bindings[k]=b
  }
  origins[strings.TrimRight(g.Origin,"/")]=true
 }
 result:=make([]string,0,len(origins));for origin:=range origins{result=append(result,origin)};sort.Strings(result)
 // The current readiness handler has a bounded aggregate probe budget.
 if len(result)>1{return nil,nil,errors.New("school grants must share one Platform origin per deployment")}
 return bindings,result,nil
}

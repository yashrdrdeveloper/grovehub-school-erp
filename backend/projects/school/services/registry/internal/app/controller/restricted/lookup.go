package restricted

import (
 "context"
 "database/sql"
 "encoding/json"
 "errors"
 "io"
 "net/http"
 "strings"
 "time"

 core "grovehub/libs/core/gateway"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/registry/internal/app/repository"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/registry/internal/app/services"
)

const Pattern="POST /school/restricted/registry/lookup"
const Scope="internal.school.registry.read"
const Subject="internal-tool-school-directory"

type Grant struct {
 Key core.WorkloadKeyConfig `json:"key"`
 SchoolIDs []string `json:"schoolIds"`
}

// New constructs independent credentials and a per-credential school allowlist.
func New(db *sql.DB,raw string)(http.HandlerFunc,error){
 if strings.TrimSpace(raw)==""{return func(w http.ResponseWriter,_ *http.Request){w.WriteHeader(403)},nil}
 if len(raw)>256*1024{return nil,errors.New("invalid school registry grants")}
 var grants []Grant;var extra any
 d:=json.NewDecoder(strings.NewReader(raw));d.DisallowUnknownFields()
 if d.Decode(&grants)!=nil||d.Decode(&extra)!=io.EOF||len(grants)==0||len(grants)>100{return nil,errors.New("invalid school registry grants")}
 keys:=[]core.WorkloadKeyConfig{};allowed:=map[string]map[string]bool{};public:=map[string]bool{}
 for _,g:=range grants{
  k:=g.Key
  if k.Issuer!=Subject||!strings.HasPrefix(k.Subject,Subject+":")||!services.ValidID(strings.TrimPrefix(k.Subject,Subject+":"))||k.AllowDelegation||len(k.Scopes)!=1||k.Scopes[0]!=Scope||public[k.PublicKey]||allowed[k.Subject]!=nil||len(g.SchoolIDs)==0||len(g.SchoolIDs)>100{return nil,errors.New("invalid school registry grant")}
  public[k.PublicKey]=true;keys=append(keys,k)
  ids:=map[string]bool{}
  for _,id:=range g.SchoolIDs{if !services.ValidID(id)||ids[id]{return nil,errors.New("invalid school registry grant IDs")};ids[id]=true}
  allowed[k.Subject]=ids
 }
 credentials,err:=core.NewEd25519Credentials(keys);if err!=nil{return nil,err}
 verifier:=core.ScopedWorkloadVerifier{Credentials:credentials,Audience:"school",Routes:map[string]core.RoutePolicy{Pattern:{Scope:Scope}}}
 return func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Cache-Control","no-store")
  if r.Method!="POST"||r.URL.Path!="/school/restricted/registry/lookup"||r.URL.RawPath!=""||r.URL.RawQuery!=""||r.URL.ForceQuery||len(r.Header.Values("Cookie"))!=0||len(r.Header.Values("X-GH-Identity-Proof"))!=0{w.WriteHeader(403);return}
  ctx,err:=verifier.VerifyWorkload(r.Context(),r);if err!=nil{w.WriteHeader(403);return}
  claims,ok:=core.WorkloadFrom(ctx)
  if !ok||len(claims.Scopes)!=1||claims.Scopes[0]!=Scope||claims.ActorGID!=""||claims.TenantID!=""||claims.SessionID!=""{w.WriteHeader(403);return}
  var in struct{SchoolIDs []string `json:"schoolIds"`};var extra any
  d:=json.NewDecoder(http.MaxBytesReader(w,r.Body,8192));d.DisallowUnknownFields()
  if d.Decode(&in)!=nil||d.Decode(&extra)!=io.EOF||len(in.SchoolIDs)<1||len(in.SchoolIDs)>100{w.WriteHeader(400);return}
  seen:=map[string]bool{}
  for _,id:=range in.SchoolIDs{
   if !services.ValidID(id)||seen[id]{w.WriteHeader(400);return};seen[id]=true
   if !allowed[claims.Subject][id]{w.WriteHeader(403);return}
  }
  probe,cancel:=context.WithTimeout(ctx,3*time.Second);defer cancel()
  rows,err:=repository.Lookup(probe,db,in.SchoolIDs);if err!=nil{w.WriteHeader(503);return}
  w.Header().Set("Content-Type","application/json")
  _=json.NewEncoder(w).Encode(map[string]any{"records":rows})
 },nil
}

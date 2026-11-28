package restricted

import (
 "context"
 "database/sql"
 "encoding/json"
 "errors"
 "io"
 "net/http"
 "strconv"
 "strings"
 "time"

 core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/registry/internal/app/services"
)

const AccessPattern="GET /school/restricted/registry/schools/{schoolId}/access-check"
const AccessScope="platform.school.access-check"
const AccessSubject="platform-tenant-access"

func NewAccess(db *sql.DB,raw,readRaw,renameRaw,statusRaw string)(http.HandlerFunc,error){
 if strings.TrimSpace(raw)==""{return func(w http.ResponseWriter,_ *http.Request){w.WriteHeader(403)},nil}
 if len(raw)>65536{return nil,errors.New("invalid access keys")}
 var keys []core.WorkloadKeyConfig;var extra any
 d:=json.NewDecoder(strings.NewReader(raw));d.DisallowUnknownFields()
 if d.Decode(&keys)!=nil||d.Decode(&extra)!=io.EOF||len(keys)==0||len(keys)>16{return nil,errors.New("invalid access keys")}
 public:=map[string]bool{};ids:=map[string]bool{}
 reserve:=func(k core.WorkloadKeyConfig){public[k.PublicKey]=true;ids[k.KeyID]=true}
 var readers []Grant
 if strings.TrimSpace(readRaw)!=""&&json.Unmarshal([]byte(readRaw),&readers)!=nil{return nil,errors.New("invalid reader keys")}
 for _,r:=range readers{reserve(r.Key)}
 for _,value:=range []string{renameRaw,statusRaw}{
  var other []core.WorkloadKeyConfig
  if strings.TrimSpace(value)!=""&&json.Unmarshal([]byte(value),&other)!=nil{return nil,errors.New("invalid mutation keys")}
  for _,k:=range other{reserve(k)}
 }
 for _,k:=range keys{
  if k.Issuer!=AccessSubject||k.Subject!=AccessSubject||k.AllowDelegation||len(k.Scopes)!=1||k.Scopes[0]!=AccessScope||public[k.PublicKey]||ids[k.KeyID]{return nil,errors.New("invalid or reused access credential")}
  reserve(k)
 }
 credentials,err:=core.NewEd25519Credentials(keys);if err!=nil{return nil,err}
 verifier:=core.ScopedWorkloadVerifier{Credentials:credentials,Audience:"school",Routes:map[string]core.RoutePolicy{AccessPattern:{Scope:AccessScope}}}
 return func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Cache-Control","no-store")
  id:=r.PathValue("schoolId")
  if !services.ValidID(id)||r.Method!="GET"||r.URL.Path!="/school/restricted/registry/schools/"+id+"/access-check"||r.URL.RawPath!=""||r.URL.RawQuery!=""||r.URL.ForceQuery||r.ContentLength!=0||len(r.TransferEncoding)!=0||(r.Body!=nil&&r.Body!=http.NoBody)||len(r.Header.Values("Cookie"))!=0||len(r.Header.Values("X-GH-Identity-Proof"))!=0{w.WriteHeader(403);return}
  ctx,err:=verifier.VerifyWorkload(r.Context(),r);if err!=nil{w.WriteHeader(403);return}
  claims,ok:=core.WorkloadFrom(ctx)
  if !ok||claims.Subject!=AccessSubject||len(claims.Scopes)!=1||claims.Scopes[0]!=AccessScope||claims.ActorGID!=""||claims.TenantID!=""||claims.SessionID!=""{w.WriteHeader(403);return}
  probe,cancel:=context.WithTimeout(ctx,2*time.Second);defer cancel()
  var generation int64
  err=db.QueryRowContext(probe,"SELECT session_generation FROM school_registry.schools WHERE id=$1 AND status='ACTIVE'",id).Scan(&generation)
  if errors.Is(err,sql.ErrNoRows){w.WriteHeader(403);return}
  if err!=nil||generation<1{w.WriteHeader(503);return}
  w.Header().Set("X-GH-Session-Generation",strconv.FormatInt(generation,10))
  w.WriteHeader(204)
 },nil
}

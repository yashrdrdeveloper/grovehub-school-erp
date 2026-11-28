package schoolexecution

import (
 "grovehub/launcher/internal/src/services/gateway/operations/workload/helpers"
 "context"
 "errors"
 "net/http"
 "strings"

 core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
)

const SchoolAuditCheckRoute="GET /audit/restricted/schools/{schoolId}/activity-check"
const SchoolAuditCheckScope="internal.audit.school.activity-check"
const SchoolExecutionIssuer="internal-tool-school"
const SchoolSubjectPrefix=SchoolExecutionIssuer+":"


func New(keys,other []core.WorkloadKeyConfig)(core.WorkloadVerifier,error){
 for i,key:=range keys{
  if key.Issuer!=SchoolExecutionIssuer||!strings.HasPrefix(key.Subject,SchoolSubjectPrefix)||!helpers.ValidSchoolID(strings.TrimPrefix(key.Subject,SchoolSubjectPrefix))||key.AllowDelegation||len(key.Scopes)!=1||key.Scopes[0]!=SchoolAuditCheckScope{return nil,errors.New("invalid school execution key")}
  for _,k:=range append(append([]core.WorkloadKeyConfig(nil),other...),keys[:i]...){
   if k.KeyID==key.KeyID||k.PublicKey==key.PublicKey{return nil,errors.New("school execution key must be independent")}
  }
 }
 credentials,err:=core.NewEd25519Credentials(keys);if err!=nil{return nil,err}
 return &verifier{scoped:core.ScopedWorkloadVerifier{Credentials:credentials,Audience:"platform",Routes:map[string]core.RoutePolicy{SchoolAuditCheckRoute:{Scope:SchoolAuditCheckScope}}}},nil
}
func(v *verifier) VerifyWorkload(ctx context.Context,r *http.Request)(context.Context,error){
 if r.URL==nil{return nil,core.ErrWorkloadDenied}
 // ServeMux.Handler resolves the route pattern but does not populate path values.
 school:=strings.TrimSuffix(strings.TrimPrefix(r.URL.Path,"/audit/restricted/schools/"),"/activity-check")
 if !helpers.ValidSchoolID(school)||r.Method!="GET"||r.URL.Path!="/audit/restricted/schools/"+school+"/activity-check"||r.URL.RawPath!=""||r.URL.RawQuery!=""||r.URL.ForceQuery||r.ContentLength!=0||len(r.TransferEncoding)!=0||(r.Body!=nil&&r.Body!=http.NoBody)||len(r.Header.Values("Cookie"))!=0||len(r.Header.Values("X-GH-Identity-Proof"))!=0{return nil,core.ErrWorkloadDenied}
 verified,err:=v.scoped.VerifyWorkload(ctx,r);if err!=nil{return nil,err}
 claims,ok:=core.WorkloadFrom(verified)
 if !ok||claims.Subject!=SchoolSubjectPrefix+school||len(claims.Scopes)!=1||claims.Scopes[0]!=SchoolAuditCheckScope||claims.ActorGID!=""||claims.TenantID!=""||claims.SessionID!=""{return nil,core.ErrWorkloadDenied}
 return verified,nil
}

// @Strategy
type verifier struct { scoped core.ScopedWorkloadVerifier }
var _ core.WorkloadVerifier = (*verifier)(nil)

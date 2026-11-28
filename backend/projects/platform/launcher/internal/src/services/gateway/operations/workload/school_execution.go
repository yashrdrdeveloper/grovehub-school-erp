package workload

import (
 "context"
 "errors"
 "net/http"
 "strings"

 "github.com/google/uuid"
 core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
)

const SchoolAuditCheckRoute="GET /audit/restricted/schools/{schoolId}/activity-check"
const SchoolAuditCheckScope="internal.audit.school.activity-check"
const SchoolExecutionIssuer="internal-tool-school"
const SchoolSubjectPrefix=SchoolExecutionIssuer+":"

func validSchoolID(value string)bool{
 id,err:=uuid.Parse(value)
 return err==nil&&id!=uuid.Nil&&id.String()==value
}

func NewSchoolExecutionVerifier(keys,other []core.WorkloadKeyConfig)(core.ScopedWorkloadVerifier,error){
 for i,key:=range keys{
  if key.Issuer!=SchoolExecutionIssuer||!strings.HasPrefix(key.Subject,SchoolSubjectPrefix)||!validSchoolID(strings.TrimPrefix(key.Subject,SchoolSubjectPrefix))||key.AllowDelegation||len(key.Scopes)!=1||key.Scopes[0]!=SchoolAuditCheckScope{return core.ScopedWorkloadVerifier{},errors.New("invalid school execution key")}
  for _,k:=range append(append([]core.WorkloadKeyConfig(nil),other...),keys[:i]...){
   if k.KeyID==key.KeyID||k.PublicKey==key.PublicKey{return core.ScopedWorkloadVerifier{},errors.New("school execution key must be independent")}
  }
 }
 credentials,err:=core.NewEd25519Credentials(keys);if err!=nil{return core.ScopedWorkloadVerifier{},err}
 return core.ScopedWorkloadVerifier{Credentials:credentials,Audience:"platform",Routes:map[string]core.RoutePolicy{SchoolAuditCheckRoute:{Scope:SchoolAuditCheckScope}}},nil
}
func(v *platformWorkloads) verifySchoolExecution(ctx context.Context,r *http.Request)(context.Context,error){
 if r.URL==nil{return nil,core.ErrWorkloadDenied}
 // ServeMux.Handler resolves the route pattern but does not populate path values.
 school:=strings.TrimSuffix(strings.TrimPrefix(r.URL.Path,"/audit/restricted/schools/"),"/activity-check")
 if !validSchoolID(school)||r.Method!="GET"||r.URL.Path!="/audit/restricted/schools/"+school+"/activity-check"||r.URL.RawPath!=""||r.URL.RawQuery!=""||r.URL.ForceQuery||r.ContentLength!=0||len(r.TransferEncoding)!=0||(r.Body!=nil&&r.Body!=http.NoBody)||len(r.Header.Values("Cookie"))!=0||len(r.Header.Values("X-GH-Identity-Proof"))!=0{return nil,core.ErrWorkloadDenied}
 verified,err:=v.schoolExecution.VerifyWorkload(ctx,r);if err!=nil{return nil,err}
 claims,ok:=core.WorkloadFrom(verified)
 if !ok||claims.Subject!=SchoolSubjectPrefix+school||len(claims.Scopes)!=1||claims.Scopes[0]!=SchoolAuditCheckScope||claims.ActorGID!=""||claims.TenantID!=""||claims.SessionID!=""{return nil,core.ErrWorkloadDenied}
 return verified,nil
}

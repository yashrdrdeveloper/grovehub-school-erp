package gateway

import (
 "context"
 "errors"
 "net/http"
 "strings"

 "github.com/google/uuid"
 core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
)

const schoolAuditCheckRoute="GET /audit/restricted/schools/{schoolId}/activity-check"
const schoolAuditCheckScope="internal.audit.school.activity-check"
const schoolExecutionIssuer="internal-tool-school"
const schoolSubjectPrefix=schoolExecutionIssuer+":"

func validSchoolID(value string)bool{
 id,err:=uuid.Parse(value)
 return err==nil&&id!=uuid.Nil&&id.String()==value
}

func newSchoolExecutionVerifier(keys,other []core.WorkloadKeyConfig)(core.ScopedWorkloadVerifier,error){
 for i,key:=range keys{
  if key.Issuer!=schoolExecutionIssuer||!strings.HasPrefix(key.Subject,schoolSubjectPrefix)||!validSchoolID(strings.TrimPrefix(key.Subject,schoolSubjectPrefix))||key.AllowDelegation||len(key.Scopes)!=1||key.Scopes[0]!=schoolAuditCheckScope{return core.ScopedWorkloadVerifier{},errors.New("invalid school execution key")}
  for _,k:=range append(append([]core.WorkloadKeyConfig(nil),other...),keys[:i]...){
   if k.KeyID==key.KeyID||k.PublicKey==key.PublicKey{return core.ScopedWorkloadVerifier{},errors.New("school execution key must be independent")}
  }
 }
 credentials,err:=core.NewEd25519Credentials(keys);if err!=nil{return core.ScopedWorkloadVerifier{},err}
 return core.ScopedWorkloadVerifier{Credentials:credentials,Audience:"platform",Routes:map[string]core.RoutePolicy{schoolAuditCheckRoute:{Scope:schoolAuditCheckScope}}},nil
}
func(v *platformWorkloads) verifySchoolExecution(ctx context.Context,r *http.Request)(context.Context,error){
 school:=r.PathValue("schoolId")
 if !validSchoolID(school)||r.Method!="GET"||r.URL==nil||r.URL.Path!="/audit/restricted/schools/"+school+"/activity-check"||r.URL.RawPath!=""||r.URL.RawQuery!=""||r.URL.ForceQuery||r.ContentLength!=0||len(r.Header.Values("Cookie"))!=0||len(r.Header.Values("X-GH-Identity-Proof"))!=0{return nil,core.ErrWorkloadDenied}
 verified,err:=v.schoolExecution.VerifyWorkload(ctx,r);if err!=nil{return nil,err}
 claims,ok:=core.WorkloadFrom(verified)
 if !ok||claims.Subject!=schoolSubjectPrefix+school||len(claims.Scopes)!=1||claims.Scopes[0]!=schoolAuditCheckScope||claims.ActorGID!=""||claims.TenantID!=""||claims.SessionID!=""{return nil,core.ErrWorkloadDenied}
 return verified,nil
}

package workload

import (
 "context"
 "errors"
 "net/http"
 "strings"
 core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
)
const StatusReceiptRoute="POST /audit/restricted/school-status-receipts"
const StatusReceiptScope="internal.audit.school.status-receipt.write"
const StatusReceiptIssuer="school-registry-status-audit"
func NewStatusReceiptVerifier(keys,other []core.WorkloadKeyConfig)(core.ScopedWorkloadVerifier,error){
 for i,k:=range keys{
  if k.Issuer!=StatusReceiptIssuer||!strings.HasPrefix(k.Subject,StatusReceiptIssuer+":")||!validSchoolID(strings.TrimPrefix(k.Subject,StatusReceiptIssuer+":"))||k.AllowDelegation||len(k.Scopes)!=1||k.Scopes[0]!=StatusReceiptScope{return core.ScopedWorkloadVerifier{},errors.New("invalid school receipt key")}
  for _,old:=range append(append([]core.WorkloadKeyConfig(nil),other...),keys[:i]...){if old.KeyID==k.KeyID||old.PublicKey==k.PublicKey{return core.ScopedWorkloadVerifier{},errors.New("school receipt key must be independent")}}
 }
 credentials,err:=core.NewEd25519Credentials(keys);if err!=nil{return core.ScopedWorkloadVerifier{},err}
 return core.ScopedWorkloadVerifier{Credentials:credentials,Audience:"platform",Routes:map[string]core.RoutePolicy{StatusReceiptRoute:{Scope:StatusReceiptScope}}},nil
}
func(v *platformWorkloads) verifyStatusReceipt(ctx context.Context,r *http.Request)(context.Context,error){
 if r.Method!="POST"||r.URL==nil||r.URL.Path!="/audit/restricted/school-status-receipts"||r.URL.RawPath!=""||r.URL.RawQuery!=""||r.URL.ForceQuery||len(r.Header.Values("Cookie"))!=0||len(r.Header.Values("X-GH-Identity-Proof"))!=0{return nil,core.ErrWorkloadDenied}
 verified,err:=v.schoolStatusReceipt.VerifyWorkload(ctx,r);if err!=nil{return nil,err}
 claims,ok:=core.WorkloadFrom(verified)
 if !ok||!strings.HasPrefix(claims.Subject,StatusReceiptIssuer+":")||!validSchoolID(strings.TrimPrefix(claims.Subject,StatusReceiptIssuer+":"))||len(claims.Scopes)!=1||claims.Scopes[0]!=StatusReceiptScope||claims.ActorGID!=""||claims.TenantID!=""||claims.SessionID!=""{return nil,core.ErrWorkloadDenied}
 return verified,nil
}

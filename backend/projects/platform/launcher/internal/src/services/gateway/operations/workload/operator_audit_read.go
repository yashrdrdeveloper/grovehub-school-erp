package workload

import (
 "context"
 "errors"
 "net/http"
 core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
)
const OperatorAuditReadRoute="POST /audit/restricted/operator-activities/search"
const OperatorAuditReadScope="internal.audit.operator.read"
func NewOperatorAuditReadVerifier(keys,other []core.WorkloadKeyConfig)(core.ScopedWorkloadVerifier,error){
 for _,key:=range keys{
  if key.Subject!="internal-tool-audit-reader"||key.Issuer!="internal-tool-audit-reader"||key.AllowDelegation||len(key.Scopes)!=1||key.Scopes[0]!=OperatorAuditReadScope{return core.ScopedWorkloadVerifier{},errors.New("invalid operator audit reader key")}
  for _,k:=range other{if k.KeyID==key.KeyID||k.PublicKey==key.PublicKey{return core.ScopedWorkloadVerifier{},errors.New("operator audit reader key must be independent")}}
 }
 credentials,err:=core.NewEd25519Credentials(keys);if err!=nil{return core.ScopedWorkloadVerifier{},err}
 return core.ScopedWorkloadVerifier{Credentials:credentials,Audience:"platform",Routes:map[string]core.RoutePolicy{OperatorAuditReadRoute:{Scope:OperatorAuditReadScope}}},nil
}
func(v *platformWorkloads) verifyOperatorAuditRead(ctx context.Context,r *http.Request)(context.Context,error){
 if r.Method!="POST"||r.URL==nil||r.URL.Path!="/audit/restricted/operator-activities/search"||r.URL.RawQuery!=""||r.URL.ForceQuery||r.URL.RawPath!=""||len(r.Header.Values("Cookie"))!=0||len(r.Header.Values("X-GH-Identity-Proof"))!=0{return nil,core.ErrWorkloadDenied}
 verified,err:=v.operatorAuditRead.VerifyWorkload(ctx,r);if err!=nil{return nil,err}
 claims,ok:=core.WorkloadFrom(verified)
 if !ok||claims.Subject!="internal-tool-audit-reader"||len(claims.Scopes)!=1||claims.Scopes[0]!=OperatorAuditReadScope||claims.ActorGID!=""||claims.TenantID!=""||claims.SessionID!=""{return nil,core.ErrWorkloadDenied}
 return verified,nil
}

package gateway

import (
 "context"
 "errors"
 "net/http"
 "strings"
 core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
)
const statusReceiptRoute="POST /audit/restricted/school-status-receipts"
const statusReceiptScope="internal.audit.school.status-receipt.write"
const statusReceiptIssuer="school-registry-status-audit"
func newStatusReceiptVerifier(keys,other []core.WorkloadKeyConfig)(core.ScopedWorkloadVerifier,error){
 for i,k:=range keys{
  if k.Issuer!=statusReceiptIssuer||!strings.HasPrefix(k.Subject,statusReceiptIssuer+":")||!validSchoolID(strings.TrimPrefix(k.Subject,statusReceiptIssuer+":"))||k.AllowDelegation||len(k.Scopes)!=1||k.Scopes[0]!=statusReceiptScope{return core.ScopedWorkloadVerifier{},errors.New("invalid school receipt key")}
  for _,old:=range append(append([]core.WorkloadKeyConfig(nil),other...),keys[:i]...){if old.KeyID==k.KeyID||old.PublicKey==k.PublicKey{return core.ScopedWorkloadVerifier{},errors.New("school receipt key must be independent")}}
 }
 credentials,err:=core.NewEd25519Credentials(keys);if err!=nil{return core.ScopedWorkloadVerifier{},err}
 return core.ScopedWorkloadVerifier{Credentials:credentials,Audience:"platform",Routes:map[string]core.RoutePolicy{statusReceiptRoute:{Scope:statusReceiptScope}}},nil
}
func(v *platformWorkloads) verifyStatusReceipt(ctx context.Context,r *http.Request)(context.Context,error){
 if r.Method!="POST"||r.URL==nil||r.URL.Path!="/audit/restricted/school-status-receipts"||r.URL.RawPath!=""||r.URL.RawQuery!=""||r.URL.ForceQuery||len(r.Header.Values("Cookie"))!=0||len(r.Header.Values("X-GH-Identity-Proof"))!=0{return nil,core.ErrWorkloadDenied}
 verified,err:=v.schoolStatusReceipt.VerifyWorkload(ctx,r);if err!=nil{return nil,err}
 claims,ok:=core.WorkloadFrom(verified)
 if !ok||!strings.HasPrefix(claims.Subject,statusReceiptIssuer+":")||!validSchoolID(strings.TrimPrefix(claims.Subject,statusReceiptIssuer+":"))||len(claims.Scopes)!=1||claims.Scopes[0]!=statusReceiptScope||claims.ActorGID!=""||claims.TenantID!=""||claims.SessionID!=""{return nil,core.ErrWorkloadDenied}
 return verified,nil
}

package gateway

import (
 "context"
 "errors"
 "net/http"

 core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
)

const internalCatalogRoute="GET /integrations/restricted/providers"
const internalCatalogScope="internal.integration.catalog.read"

func newInternalCatalogVerifier(keys,schoolKeys []core.WorkloadKeyConfig)(core.ScopedWorkloadVerifier,error){
 for _,key:=range keys{
  if key.Subject!="internal-tool"||key.Issuer!="internal-tool"||key.AllowDelegation||len(key.Scopes)!=1||key.Scopes[0]!=internalCatalogScope{return core.ScopedWorkloadVerifier{},errors.New("invalid internal tool catalog key")}
  for _,school:=range schoolKeys{if school.KeyID==key.KeyID||school.PublicKey==key.PublicKey{return core.ScopedWorkloadVerifier{},errors.New("internal tool key must be independent")}}
 }
 credentials,err:=core.NewEd25519Credentials(keys)
 if err!=nil{return core.ScopedWorkloadVerifier{},err}
 return core.ScopedWorkloadVerifier{Credentials:credentials,Audience:"platform",Routes:map[string]core.RoutePolicy{internalCatalogRoute:{Scope:internalCatalogScope}}},nil
}

func(v *platformWorkloads) verifyInternalCatalog(ctx context.Context,r *http.Request)(context.Context,error){
 if r.Pattern!=internalCatalogRoute||r.Method!="GET"||r.URL==nil||r.URL.Path!="/integrations/restricted/providers"||r.URL.RawQuery!=""||r.URL.RawPath!=""||len(r.Header.Values("Cookie"))!=0||len(r.Header.Values("X-GH-Identity-Proof"))!=0{return nil,core.ErrWorkloadDenied}
 verified,err:=v.internalCatalog.VerifyWorkload(ctx,r);if err!=nil{return nil,err}
 claims,ok:=core.WorkloadFrom(verified)
 if !ok||claims.Subject!="internal-tool"||len(claims.Scopes)!=1||claims.Scopes[0]!=internalCatalogScope||claims.ActorGID!=""||claims.TenantID!=""||claims.SessionID!=""{return nil,core.ErrWorkloadDenied}
 return verified,nil
}

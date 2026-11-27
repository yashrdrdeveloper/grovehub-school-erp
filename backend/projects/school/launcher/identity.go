package launcher

import (
 "crypto/ed25519"
 "encoding/base64"
 "fmt"
 "net/http"
 "net/url"
 "os"

 core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
 client "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/identity/client/authentication"
 controller "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/identity/controller/authentication/public"
 protected "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/identity/controller/authentication/protected"
 "strings"
 business "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/platform/client"
 businesscontroller "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/platform/controller/protected"
 audit "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/audit/client/activity"
 integration "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/integration/client/connection"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/internal/platformhttp"
)

func registerIdentity(mux *http.ServeMux) error {
 origin,err:=url.Parse(os.Getenv("PLATFORM_ORIGIN"))
 if err!=nil || origin.Host=="" || origin.User!=nil || (origin.Path!="" && origin.Path!="/") || origin.RawQuery!="" || origin.Fragment!="" {return fmt.Errorf("invalid Platform origin")}
 browserOrigin:=os.Getenv("SCHOOL_ALLOWED_ORIGIN")
 if browserOrigin=="" {return fmt.Errorf("School origin required")}
 key,err:=base64.RawURLEncoding.DecodeString(os.Getenv("SCHOOL_WORKLOAD_PRIVATE_KEY"))
 if err!=nil || len(key)!=ed25519.PrivateKeySize {return fmt.Errorf("invalid School workload key")}
 transport:=&core.WorkloadTransport{Origin:origin,KeyID:os.Getenv("SCHOOL_WORKLOAD_KEY_ID"),Issuer:"school",Subject:"school",Audience:"platform",PrivateKey:ed25519.PrivateKey(key),Base:platformhttp.ProofTransport{},Policy:func(r *http.Request)(core.WorkloadClaims,error){
  switch r.Method+" "+r.URL.Path {
  case "POST /identity/restricted/auth/login": return core.WorkloadClaims{Scopes:[]string{"identity.auth.login"}},nil
  case "POST /identity/restricted/auth/forgot-password": return core.WorkloadClaims{Scopes:[]string{"identity.auth.forgot"}},nil
  case "POST /identity/restricted/auth/refresh": return core.WorkloadClaims{Scopes:[]string{"identity.auth.refresh"}},nil
  case "GET /identity/restricted/auth/recovery/validate": return core.WorkloadClaims{Scopes:[]string{"identity.auth.recovery.validate"}},nil
  case "GET /identity/restricted/auth/invitations/validate": return core.WorkloadClaims{Scopes:[]string{"identity.auth.invitation.validate"}},nil
  case "POST /identity/restricted/auth/recovery/exchange": return core.WorkloadClaims{Scopes:[]string{"identity.auth.recovery.exchange"}},nil
  case "POST /identity/restricted/auth/invitations/exchange": return core.WorkloadClaims{Scopes:[]string{"identity.auth.invitation.exchange"}},nil
  case "GET /identity/restricted/auth/reset-password/status": return core.WorkloadClaims{Scopes:[]string{"identity.auth.reset.status"}},nil
  case "POST /identity/restricted/auth/reset-password": return core.WorkloadClaims{Scopes:[]string{"identity.auth.reset"}},nil
  case "GET /identity/restricted/accounts/setup/status": return core.WorkloadClaims{Scopes:[]string{"identity.auth.setup.status"}},nil
  case "POST /identity/restricted/accounts/setup": return core.WorkloadClaims{Scopes:[]string{"identity.auth.setup"}},nil
  case "POST /identity/restricted/auth/logout": return core.WorkloadClaims{Scopes:[]string{"identity.auth.logout"}},nil
  case "POST /identity/restricted/auth/logout-all": return core.WorkloadClaims{Scopes:[]string{"identity.auth.logout.all"}},nil
  case "POST /identity/restricted/auth/change-password": return core.WorkloadClaims{Scopes:[]string{"identity.auth.change.password"}},nil
  case "POST /identity/restricted/auth/reauthenticate": return core.WorkloadClaims{Scopes:[]string{"identity.auth.reauthenticate"}},nil
  case "GET /identity/restricted/auth/sessions": return core.WorkloadClaims{Scopes:[]string{"identity.auth.sessions"}},nil
  case "GET /identity/restricted/auth/me": return core.WorkloadClaims{Scopes:[]string{"identity.auth.me"}},nil
  case "GET /integrations/restricted/connections/oauth/gmail/callback", "GET /integrations/restricted/connections/oauth/outlook/callback": return core.WorkloadClaims{Scopes:[]string{"integration.oauth.callback"}},nil
  default:
   if r.Method=="DELETE" && strings.HasPrefix(r.URL.Path,"/identity/restricted/auth/sessions/") && client.ValidSessionID(strings.TrimPrefix(r.URL.Path,"/identity/restricted/auth/sessions/")) {return core.WorkloadClaims{Scopes:[]string{"identity.auth.sessions.revoke"}},nil}
   if scope,ok:=business.Scope(r);ok {return core.WorkloadClaims{Scopes:[]string{scope}},nil}
   return core.WorkloadClaims{},core.ErrWorkloadDenied
  }
 }}
 httpClient,err:=core.NewWorkloadHTTPClient(transport)
 if err!=nil {return err}
 platform:=client.New(origin,httpClient)
 controller.New(platform,browserOrigin).Register(mux)
 protected.New(platform,browserOrigin).Register(mux)
 connections:=integration.New(origin,httpClient)
 businesscontroller.Register(mux,platform,audit.New(origin,httpClient),connections,browserOrigin)
 businesscontroller.RegisterOAuthCallback(mux,connections)
 return nil
}

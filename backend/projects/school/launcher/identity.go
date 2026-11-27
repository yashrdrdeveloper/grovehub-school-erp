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
)

func registerIdentity(mux *http.ServeMux) error {
 origin,err:=url.Parse(os.Getenv("PLATFORM_ORIGIN"))
 if err!=nil || origin.Host=="" || origin.User!=nil || (origin.Path!="" && origin.Path!="/") || origin.RawQuery!="" || origin.Fragment!="" {return fmt.Errorf("invalid Platform origin")}
 browserOrigin:=os.Getenv("SCHOOL_ALLOWED_ORIGIN")
 if browserOrigin=="" {return fmt.Errorf("School origin required")}
 key,err:=base64.RawURLEncoding.DecodeString(os.Getenv("SCHOOL_WORKLOAD_PRIVATE_KEY"))
 if err!=nil || len(key)!=ed25519.PrivateKeySize {return fmt.Errorf("invalid School workload key")}
 transport:=&core.WorkloadTransport{Origin:origin,KeyID:os.Getenv("SCHOOL_WORKLOAD_KEY_ID"),Issuer:"school",Subject:"school",Audience:"platform",PrivateKey:ed25519.PrivateKey(key),Base:client.ProofTransport{},Policy:func(r *http.Request)(core.WorkloadClaims,error){
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
  default: return core.WorkloadClaims{},core.ErrWorkloadDenied
  }
 }}
 httpClient,err:=core.NewWorkloadHTTPClient(transport)
 if err!=nil {return err}
 controller.New(client.New(origin,httpClient),browserOrigin).Register(mux)
 return nil
}

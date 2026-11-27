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
 transport:=&core.WorkloadTransport{Origin:origin,KeyID:os.Getenv("SCHOOL_WORKLOAD_KEY_ID"),Issuer:"school",Subject:"school",Audience:"platform",PrivateKey:ed25519.PrivateKey(key),Policy:func(r *http.Request)(core.WorkloadClaims,error){
  switch r.Method+" "+r.URL.Path {
  case "POST /identity/restricted/auth/login": return core.WorkloadClaims{Scopes:[]string{"identity.auth.login"}},nil
  case "POST /identity/restricted/auth/forgot-password": return core.WorkloadClaims{Scopes:[]string{"identity.auth.forgot"}},nil
  default: return core.WorkloadClaims{},core.ErrWorkloadDenied
  }
 }}
 httpClient,err:=core.NewWorkloadHTTPClient(transport)
 if err!=nil {return err}
 controller.New(client.New(origin,httpClient),browserOrigin).Register(mux)
 return nil
}

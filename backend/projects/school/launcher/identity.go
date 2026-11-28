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
 newClient:=func(policy workloadPolicy)(*http.Client,error){
  return core.NewWorkloadHTTPClient(&core.WorkloadTransport{Origin:origin,KeyID:os.Getenv("SCHOOL_WORKLOAD_KEY_ID"),Issuer:"school",Subject:"school",Audience:"platform",PrivateKey:ed25519.PrivateKey(key),Base:platformhttp.ProofTransport{},Policy:policy})
 }
 identityHTTP,err:=newClient(identityPolicy)
 if err!=nil{return err}
 auditHTTP,err:=newClient(auditPolicy)
 if err!=nil{return err}
 integrationHTTP,err:=newClient(integrationPolicy)
 if err!=nil{return err}
 platform:=client.New(origin,identityHTTP)
 controller.New(platform,browserOrigin).Register(mux)
 protected.New(platform,browserOrigin).Register(mux)
 connections:=integration.New(origin,integrationHTTP)
 businesscontroller.Register(mux,platform,audit.New(origin,auditHTTP),connections,browserOrigin)
 businesscontroller.RegisterOAuthCallback(mux,connections)
 return nil
}

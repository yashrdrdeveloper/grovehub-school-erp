// Command workload-config creates local development credentials. It never
// prints private keys and refuses to overwrite an existing configuration.
package main
import (
 "crypto/ed25519"
 "crypto/rand"
 "encoding/base64"
 "encoding/json"
 "fmt"
 "log"
 "os"
 "path/filepath"
 "sort"
 business "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/platform/client"
)
type policy struct { Scope string; RequireActor bool }
type keyConfig struct { KeyID, PublicKey, Issuer, Subject string; Scopes []string; AllowDelegation bool }
type config struct { Keys []keyConfig; Routes map[string]policy }
func run(dir string)error{
 if err:=os.Mkdir(dir,0700);err!=nil{return err}
 pub,priv,err:=ed25519.GenerateKey(rand.Reader);if err!=nil{return err}
 routes:=map[string]policy{
 "GET /integrations/restricted/connections/oauth/{provider}/callback":{"integration.oauth.callback",false},
 "POST /identity/restricted/auth/login":{"identity.auth.login",false},
 "POST /identity/restricted/auth/forgot-password":{"identity.auth.forgot",false},
 "POST /identity/restricted/auth/refresh":{"identity.auth.refresh",false},
 "GET /identity/restricted/auth/recovery/validate":{"identity.auth.recovery.validate",false},
 "GET /identity/restricted/auth/invitations/validate":{"identity.auth.invitation.validate",false},
 "POST /identity/restricted/auth/recovery/exchange":{"identity.auth.recovery.exchange",false},
 "POST /identity/restricted/auth/invitations/exchange":{"identity.auth.invitation.exchange",false},
 "GET /identity/restricted/auth/reset-password/status":{"identity.auth.reset.status",false},
 "POST /identity/restricted/auth/reset-password":{"identity.auth.reset",false},
 "GET /identity/restricted/accounts/setup/status":{"identity.auth.setup.status",false},
 "POST /identity/restricted/accounts/setup":{"identity.auth.setup",false},
 "POST /identity/restricted/auth/logout":{"identity.auth.logout",true},
 "POST /identity/restricted/auth/logout-all":{"identity.auth.logout.all",true},
 "POST /identity/restricted/auth/change-password":{"identity.auth.change.password",true},
 "POST /identity/restricted/auth/reauthenticate":{"identity.auth.reauthenticate",true},
 "GET /identity/restricted/auth/sessions":{"identity.auth.sessions",true},
 "DELETE /identity/restricted/auth/sessions/{sessionId}":{"identity.auth.sessions.revoke",true},
 "GET /identity/restricted/auth/me":{"identity.auth.me",true},
 }
 for _,op:=range business.Operations(){routes[op.Method+" "+op.RestrictedPath()]=policy{op.Scope,true}}
 scopes:=make([]string,0,len(routes));for _,p:=range routes{scopes=append(scopes,p.Scope)};sort.Strings(scopes)
 const id="school-dev"
 raw,err:=json.Marshal(config{Keys:[]keyConfig{{KeyID:id,PublicKey:base64.RawURLEncoding.EncodeToString(pub),Issuer:"school",Subject:"school",Scopes:scopes}},Routes:routes});if err!=nil{return err}
 if err:=os.WriteFile(filepath.Join(dir,"platform.env"),[]byte("PLATFORM_WORKLOAD_CONFIG='"+string(raw)+"'\n"),0600);err!=nil{return err}
 school:="PLATFORM_ORIGIN=http://127.0.0.1:8080\nSCHOOL_ALLOWED_ORIGIN=http://127.0.0.1:4173\nSCHOOL_WORKLOAD_KEY_ID="+id+"\nSCHOOL_WORKLOAD_PRIVATE_KEY="+base64.RawURLEncoding.EncodeToString(priv)+"\n"
 return os.WriteFile(filepath.Join(dir,"school.env"),[]byte(school),0600)
}
func main(){if len(os.Args)!=2{log.Fatal("usage: workload-config NEW_DIRECTORY")};if err:=run(os.Args[1]);err!=nil{log.Fatal(err)};fmt.Println("Development configuration created in",os.Args[1])}

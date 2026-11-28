// Cross-project lifecycle test. School and Internal Tool run as separate
// executables; Platform uses its real HTTP gateway, controllers and databases.
package lifecycle_test

import (
 "bytes"
 "context"
 "crypto/ed25519"
 "crypto/rand"
 "encoding/base64"
 "encoding/json"
 "errors"
 "io"
 "net"
 "net/http"
 "net/http/httptest"
 "os"
 "os/exec"
 "strings"
 "sync/atomic"
 "testing"
 "time"

 "github.com/google/uuid"
 "gorm.io/driver/postgres"
 "gorm.io/gorm"
 core "grovehub/libs/core/gateway"
 coredatabase "grovehub/libs/core/database"
 gateway "grovehub/launcher"
 schoolgateway "grovehub/libs/platform/session"
 "grovehub/services/audit"
 "grovehub/services/identity"
 "grovehub/services/infrastructure"
 "grovehub/services/infrastructure/client/messaging"
 "grovehub/services/integration"
 public "grovehub/services/identity/internal/app/controller/authentication/public/impl"
 protected "grovehub/services/identity/internal/app/controller/authentication/protected/impl"
 "grovehub/services/identity/internal/app/config"
 "grovehub/services/identity/internal/app/entity"
 mapper "grovehub/services/identity/internal/app/mappers/account/impl"
 accounts "grovehub/services/identity/internal/app/repository/account/impl"
 authority "grovehub/services/identity/internal/app/repository/authority/impl"
 credentialrepo "grovehub/services/identity/internal/app/repository/credential/impl"
 mailrepo "grovehub/services/identity/internal/app/repository/invitation/impl"
 sessions "grovehub/services/identity/internal/app/repository/session/impl"
 tokens "grovehub/services/identity/internal/app/repository/passwordtoken/impl"
 authimpl "grovehub/services/identity/internal/app/services/authentication/impl"
 secrets "grovehub/services/identity/internal/app/services/authentication/helpers"
 reconcile "grovehub/services/identity/internal/app/services/authentication/operations/reconciletenant"
 credentials "grovehub/services/identity/internal/app/services/credential/impl"
 "grovehub/services/identity/internal/app/services/registration/operations/invitation"
 security "grovehub/services/identity/internal/app/services/security"
 access "grovehub/services/identity/internal/app/services/security/impl"
 "grovehub/services/identity/internal/app/services/security/operations/routepolicy"
 cookieimpl "grovehub/services/identity/internal/app/transport/cookies/impl"
)

type idleBroker struct{messaging.MessagingClient}
func(idleBroker) Consume(ctx context.Context,_ string,_ messaging.MessageHandler)error{<-ctx.Done();return ctx.Err()}
func(idleBroker) Health(context.Context)error{return nil}

type identityBoundary struct{security.AccessVerifier;policy *routepolicy.Policy}
func(b identityBoundary) VerifyDelegatedSession(context.Context,string)(schoolgateway.Principal,error){return schoolgateway.Principal{},errors.New("delegation not enabled in fixture")}
func(b identityBoundary) AuthorizeProtectedRoute(ctx context.Context,pattern string)error{return b.policy.Authorize(ctx,pattern)}

type key struct{config core.WorkloadKeyConfig;private ed25519.PrivateKey}
func newKey(t *testing.T,id,issuer,subject string,scopes ...string)key{
 t.Helper();pub,priv,err:=ed25519.GenerateKey(rand.Reader);if err!=nil{t.Fatal(err)}
 return key{core.WorkloadKeyConfig{KeyID:id,PublicKey:base64.RawURLEncoding.EncodeToString(pub),Issuer:issuer,Subject:subject,Scopes:scopes},priv}
}
func(k key) secret()string{return base64.RawURLEncoding.EncodeToString(k.private)}
func jsonText(t *testing.T,v any)string{t.Helper();raw,err:=json.Marshal(v);if err!=nil{t.Fatal(err)};return string(raw)}
func must(t *testing.T,err error){t.Helper();if err!=nil{t.Fatal(err)}}
func waitFor(t *testing.T,ctx context.Context,d time.Duration,check func()bool){
 t.Helper();deadline:=time.Now().Add(d)
 for time.Now().Before(deadline){if check(){return};select{case <-ctx.Done():t.Fatal(ctx.Err());case <-time.After(200*time.Millisecond):}}
 t.Fatal("condition did not become true")
}
func runChild(t *testing.T,ctx context.Context,path string,env []string,args ...string){
 t.Helper();cmd:=exec.CommandContext(ctx,path,args...);cmd.Env=append(os.Environ(),env...)
 if out,err:=cmd.CombinedOutput();err!=nil{t.Fatalf("fixture command failed: %v\n%s",err,out)}
}

func TestCrossProjectSchoolLifecycle(t *testing.T){
 schoolBin,toolBin,dsn:=os.Getenv("LIFECYCLE_SCHOOL_BIN"),os.Getenv("LIFECYCLE_INTERNAL_TEST_BIN"),os.Getenv("LIFECYCLE_PLATFORM_DSN")
 if schoolBin==""&&toolBin==""&&dsn==""{t.Skip("requires independent disposable project fixtures")}
 if schoolBin==""||toolBin==""||dsn==""{t.Fatal("incomplete lifecycle fixture configuration")}
 ctx,cancel:=context.WithTimeout(context.Background(),3*time.Minute);defer cancel()
 db,err:=gorm.Open(postgres.Open(dsn),&gorm.Config{});must(t,err)
 sqlDB,err:=db.DB();must(t,err);defer sqlDB.Close()
 var name string
 must(t,db.Raw("SELECT current_database()").Scan(&name).Error)
 if name!="groovehub_lifecycle_test"{t.Fatal("disposable Platform database required")}
 var existing int64
 must(t,db.Raw("SELECT count(*) FROM information_schema.schemata WHERE schema_name IN ('identity','integration','audit')").Scan(&existing).Error)
 if existing!=0{t.Fatal("Platform lifecycle database must be fresh")}
 must(t,identity.Migrate(ctx,db));must(t,infrastructure.Migrate(ctx,db));must(t,integration.Migrate(ctx,db));must(t,audit.Migrate(ctx,db))

 school,other,readerID:=uuid.NewString(),uuid.NewString(),uuid.NewString()
 accountID,otherAccount:=uuid.NewString(),uuid.NewString()
 personID,err:=uuid.NewV7();must(t,err)
 otherPersonID,err:=uuid.NewV7();must(t,err)
 statusKey:=newKey(t,"lifecycle-status","internal-tool-school-status","internal-tool-school-status:"+school,"internal.school.registry.status")
 directoryKey:=newKey(t,"lifecycle-directory","internal-tool-school-directory","internal-tool-school-directory:"+readerID,"internal.school.registry.read")
 tenantKey:=newKey(t,"lifecycle-access","platform-tenant-access","platform-tenant-access","platform.school.access-check")
 receiptKey:=newKey(t,"lifecycle-receipt","school-registry-status-audit","school-registry-status-audit:"+school,"internal.audit.school.status-receipt.write")
 readerKey:=newKey(t,"lifecycle-reader","internal-tool-school-receipt-reader","internal-tool-school-receipt-reader","internal.audit.school.mutation-receipt.read")
 routes:=map[string]core.RoutePolicy{
  "POST /identity/restricted/auth/login":{Scope:"identity.auth.login"},
  "POST /identity/restricted/auth/refresh":{Scope:"identity.auth.refresh"},
  "POST /identity/restricted/auth/forgot-password":{Scope:"identity.auth.forgot"},
  "GET /identity/restricted/auth/recovery/validate":{Scope:"identity.auth.recovery.validate"},
  "POST /identity/restricted/auth/recovery/exchange":{Scope:"identity.auth.recovery.exchange"},
  "GET /identity/restricted/auth/reset-password/status":{Scope:"identity.auth.reset.status"},
  "POST /identity/restricted/auth/reset-password":{Scope:"identity.auth.reset"},
  "GET /identity/restricted/auth/me":{Scope:"identity.auth.me",RequireActor:true},
  "POST /identity/restricted/auth/logout":{Scope:"identity.auth.logout",RequireActor:true},
  "GET /integrations/restricted/connections":{Scope:"integration.list",RequireActor:true},
 }
 var scopes []string;for _,policy:=range routes{scopes=append(scopes,policy.Scope)}
 schoolKey:=newKey(t,"lifecycle-school","school","school",scopes...)
 // School's real command owns registry migrations and the HTTP listener.
 listener,err:=net.Listen("tcp","127.0.0.1:0");must(t,err)
 address:=listener.Addr().String();must(t,listener.Close())
 schoolOrigin:="http://"+address
 lookup,err:=gateway.NewTenantState(gateway.TenantAccessConfig{Origin:schoolOrigin,KeyID:tenantKey.config.KeyID,PrivateKey:tenantKey.secret()});must(t,err)
 check:=schoolgateway.ActiveTenantCheck(lookup)

 ar,sr,tr:=accounts.New(db),sessions.New(db),tokens.New(db)
 cred:=credentials.New(credentialrepo.New(db))
 const password="Lifecycle-test-password-123"
 for _,row:=range []struct{id,tenant,person string}{{accountID,school,personID.String()},{otherAccount,other,otherPersonID.String()}}{
  must(t,db.Exec("INSERT INTO identity.accounts(account_id,person_gid,tenant_id,email,status) VALUES (?,?,?,?,'ACTIVE')",row.id,row.person,row.tenant,row.id+"@example.test").Error)
  must(t,cred.Set(ctx,row.id,password))
 }
 eventStore:=infrastructure.NewEventStore(db)
 producer:=audit.NewProducer(eventStore)
 sender:=invitation.New(tr,mailrepo.New(db,producer),config.Config{EncryptionKey:make([]byte,32),ResetURL:"https://school.example/reset-password"})
 base:=authimpl.New(ar,sr,tr,cred,sender,mapper.New(),coredatabase.NewUnitOfWork(db))
 auth:=authimpl.NewRecoveryGuard(base,ar,tr,sr,authority.New(db),sender,producer)
 actor:=identityBoundary{AccessVerifier:access.New(ar,sr,nil),policy:routepolicy.NewWithExecution(nil,false,nil,check)}
 t.Setenv("GROVEHUB_INTEGRATION_ENCRYPTION_KEY",base64.StdEncoding.EncodeToString(make([]byte,32)))
 // No real provider or message broker is required for mailbox listing or
 // the synchronous receipt ledger. Consumer loops stop on cancellation.
 for _,v:=range []string{"GROVEHUB_GMAIL_CLIENT_ID","GROVEHUB_OUTLOOK_CLIENT_ID","GROVEHUB_GMAIL_PUBSUB_TOPIC","GROVEHUB_GMAIL_PUBSUB_AUDIENCE","GROVEHUB_GMAIL_PUBSUB_SERVICE_ACCOUNT"}{t.Setenv(v,"")}
 integrationModule,err:=integration.NewWithTenantState(nil,db,idleBroker{},lookup);must(t,err);defer integrationModule.Close()
 auditModule,err:=audit.New(nil,db,idleBroker{},actor);must(t,err);defer auditModule.Close()
 mux:=http.NewServeMux()
 pub:=public.New(auth,cookieimpl.New());prot:=protected.New(auth,cookieimpl.New())
 mux.HandleFunc("POST /identity/restricted/auth/login",pub.Login)
 mux.HandleFunc("POST /identity/restricted/auth/refresh",pub.Refresh)
 mux.HandleFunc("POST /identity/restricted/auth/forgot-password",pub.ForgotPassword)
 mux.HandleFunc("GET /identity/restricted/auth/recovery/validate",pub.RecoveryValidate)
 mux.HandleFunc("POST /identity/restricted/auth/recovery/exchange",pub.RecoveryExchange)
 mux.HandleFunc("GET /identity/restricted/auth/reset-password/status",pub.ResetGrantStatus)
 mux.HandleFunc("POST /identity/restricted/auth/reset-password",pub.ResetPassword)
 mux.HandleFunc("GET /identity/restricted/auth/me",prot.Me)
 mux.HandleFunc("POST /identity/restricted/auth/logout",prot.Logout)
 integrationModule.RegisterHTTP(mux);auditModule.RegisterHTTP(mux)
 mux.HandleFunc("GET /ready",func(w http.ResponseWriter,_ *http.Request){w.WriteHeader(204)})
 configJSON:=jsonText(t,map[string]any{"Keys":[]core.WorkloadKeyConfig{schoolKey.config},"Routes":routes,"SchoolStatusReceiptKeys":[]core.WorkloadKeyConfig{receiptKey.config},"SchoolReceiptReadKeys":[]core.WorkloadKeyConfig{readerKey.config}})
 workloads,err:=gateway.NewWorkloadVerifier(configJSON,actor,integrationModule,check);must(t,err)
 const browserOrigin="http://127.0.0.1:4173"
 runtime:=gateway.NewExternalGateway(mux,actor,workloads,browserOrigin)
 restricted:=core.RestrictedOnly(mux,runtime)
 var receiptAttempts atomic.Int32
 var lostAcknowledgment atomic.Bool
 platform:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.URL.Path=="/ready"{runtime.ServeHTTP(w,r);return}
  if r.URL.Path=="/audit/restricted/school-status-receipts"{
   receiptAttempts.Add(1)
   // Simulate an acknowledgment lost after the first committed receipt.
   if lostAcknowledgment.CompareAndSwap(false,true){
    rec:=httptest.NewRecorder();restricted.ServeHTTP(rec,r)
    if rec.Code==204{w.WriteHeader(503);return}
    w.WriteHeader(rec.Code);return
   }
  }
  restricted.ServeHTTP(w,r)
 }));defer platform.Close()
 schoolEnv:=[]string{
  "SCHOOL_DATABASE_URL="+os.Getenv("PROJECT_SCHOOL_TEST_DSN"),
  "GROOVEHUB_HTTP_ADDR="+address,
  "PLATFORM_ORIGIN="+platform.URL,
  "SCHOOL_ALLOWED_ORIGIN="+browserOrigin,
  "SCHOOL_WORKLOAD_KEY_ID="+schoolKey.config.KeyID,
  "SCHOOL_WORKLOAD_PRIVATE_KEY="+schoolKey.secret(),
  "SCHOOL_REGISTRY_READ_GRANTS="+jsonText(t,[]any{map[string]any{"key":directoryKey.config,"schoolIds":[]string{school,other}}}),
  "SCHOOL_REGISTRY_RENAME_KEYS=",
  "SCHOOL_REGISTRY_STATUS_KEYS="+jsonText(t,[]core.WorkloadKeyConfig{statusKey.config}),
  "SCHOOL_REGISTRY_ACCESS_KEYS="+jsonText(t,[]core.WorkloadKeyConfig{tenantKey.config}),
  "SCHOOL_RECEIPT_ORIGIN=","SCHOOL_RECEIPT_SCHOOL_ID=","SCHOOL_RECEIPT_KEY_ID=","SCHOOL_RECEIPT_PRIVATE_KEY=",
  "SCHOOL_RECEIPT_PUBLISHERS="+jsonText(t,[]any{map[string]any{"operation":"status","origin":platform.URL,"schoolId":school,"keyId":receiptKey.config.KeyID,"privateKey":receiptKey.secret()}}),
 }
 runChild(t,ctx,schoolBin,schoolEnv,"migrate")
 runChild(t,ctx,schoolBin,schoolEnv,"register-school",school,"school-"+school,"Lifecycle school")
 runChild(t,ctx,schoolBin,schoolEnv,"register-school",other,"school-"+other,"Unrelated school")
 startSchool:=func() *exec.Cmd{
  cmd:=exec.CommandContext(ctx,schoolBin);cmd.Env=append(os.Environ(),schoolEnv...)
  cmd.Stdout=os.Stdout;cmd.Stderr=os.Stderr;must(t,cmd.Start())
  t.Cleanup(func(){_ = cmd.Process.Kill();_ = cmd.Wait()})
  waitFor(t,ctx,15*time.Second,func()bool{state,err:=lookup(ctx,school);return err==nil&&state.Generation>0})
  return cmd
 }
 schoolProcess:=startSchool()
 client:=&http.Client{Timeout:10*time.Second,CheckRedirect:func(*http.Request,[]*http.Request)error{return http.ErrUseLastResponse}}
 type cookies map[string]*http.Cookie
 call:=func(method,path,body,proof string,jar cookies,want int)cookies{
  t.Helper();req,err:=http.NewRequestWithContext(ctx,method,schoolOrigin+path,strings.NewReader(body));must(t,err)
  req.Header.Set("Origin",browserOrigin);req.Header.Set("Content-Type","application/json")
  if proof!=""{req.Header.Set("Authorization","Recovery "+proof)}
  for _,cookie:=range jar{req.AddCookie(cookie)}
  if cookie:=jar["gh_csrf"];cookie!=nil{req.Header.Set("X-CSRF-Token",cookie.Value)}
  resp,err:=client.Do(req);must(t,err);defer resp.Body.Close()
  _,err=io.Copy(io.Discard,resp.Body);must(t,err)
  if resp.StatusCode!=want{t.Fatalf("%s %s: got %d want %d",method,path,resp.StatusCode,want)}
  result:=cookies{};for name,c:=range jar{result[name]=c};for _,c:=range resp.Cookies(){result[c.Name]=c};return result
 }
 login:=func(id string,want int)cookies{return call("POST","/identity/public/auth/login",jsonText(t,map[string]string{"identifier":id+"@example.test","password":password}),"",nil,want)}
 me:=func(jar cookies,want int){call("GET","/identity/protected/auth/me","","",jar,want)}
 mailbox:=func(jar cookies,want int){call("GET","/integrations/protected/connections","","",jar,want)}
 refresh:=func(jar cookies,want int){call("POST","/identity/public/auth/refresh","","",jar,want)}
 active,unrelated:=login(accountID,200),login(otherAccount,200)
 me(active,200);mailbox(active,200);me(unrelated,200)
 if _,err=integrationModule.Client().ListConnections(ctx,school,personID.String());err!=nil{t.Fatal("active direct mailbox listing failed",err)}
 // Test known proofs at generation 1 independently of asynchronous email.
 issue:=func(purpose entity.TokenPurpose,generation int64)string{
  secret,err:=secrets.NewSecret(32);must(t,err);hash,err:=secrets.Digest(secret,32);must(t,err)
  version:=int64(1);now:=time.Now().UTC()
  must(t,tr.Issue(ctx,entity.PasswordToken{TokenHash:hash,AccountID:accountID,Purpose:purpose,AccountVersion:&version,SessionGeneration:generation,CreatedAt:now,ExpiresAt:now.Add(time.Hour)}))
  return secret
 }
 proof:=issue(entity.PasswordReset,1)
 call("GET","/identity/public/auth/recovery/validate","",proof,nil,200)
 grantJar:=call("POST","/identity/public/auth/recovery/exchange","",proof,nil,200)
 if grantJar["__Secure-gh_reset"]==nil{t.Fatal("missing reset grant")}
 oldProof:=issue(entity.PasswordReset,1)
 call("GET","/identity/public/auth/reset-password/status","","",grantJar,200)
 transition:=func(expected,target string){
  t.Helper();runChild(t,ctx,toolBin,[]string{"LIFECYCLE_EXECUTE=true","LIFECYCLE_SCHOOL_ID="+school,"LIFECYCLE_SCHOOL_ORIGIN="+schoolOrigin,"LIFECYCLE_DIRECTORY_ID="+readerID,"LIFECYCLE_DIRECTORY_PRIVATE="+directoryKey.secret(),"LIFECYCLE_STATUS_PRIVATE="+statusKey.secret(),"LIFECYCLE_EXPECTED="+expected,"LIFECYCLE_TARGET="+target},"-test.run=^TestLiveSchoolLifecycleExecution$","-test.v","-test.timeout=45s")
 }
 transition("ACTIVE","DISABLED")
 state,err:=lookup(ctx,school);must(t,err);if state.Active||state.Generation!=2{t.Fatal("disable generation",state)}
 login(accountID,401);me(active,403);mailbox(active,403);refresh(active,401)
 me(unrelated,200);mailbox(unrelated,200)
 call("GET","/identity/public/auth/recovery/validate","",oldProof,nil,401)
 call("GET","/identity/public/auth/reset-password/status","","",grantJar,401)
 call("POST","/identity/public/auth/forgot-password",jsonText(t,map[string]string{"email":accountID+"@example.test"}),"",nil,202)
 var count int64;must(t,db.Raw("SELECT count(*) FROM identity.invitations WHERE account_id=?",accountID).Scan(&count).Error)
 if count!=0{t.Fatal("disabled school queued reset mail")}
 if _,err=integrationModule.Client().ListConnections(ctx,school,personID.String());err==nil{t.Fatal("direct Integration client bypassed registry")}
 transition("DISABLED","ACTIVE")
 state,err=lookup(ctx,school);must(t,err);if !state.Active||state.Generation!=2{t.Fatal("activation generation",state)}
 me(active,403);refresh(active,401);mailbox(active,403)
 call("GET","/identity/public/auth/recovery/validate","",oldProof,nil,401)
 call("GET","/identity/public/auth/reset-password/status","","",grantJar,401)
 call("POST","/identity/public/auth/reset-password",jsonText(t,map[string]string{"newPassword":"Rejected-password-456"}),"",grantJar,401)
 fresh:=login(accountID,200);me(fresh,200);mailbox(fresh,200)
 // Physical revocation is tested only after proving stale sessions fail
 // logically without the reconciliation worker.
 worker:=reconcile.New(sessions.NewTenantSessions(db),lookup,nil);must(t,worker.Step(ctx))
 must(t,db.Raw("SELECT count(*) FROM identity.sessions WHERE account_id=? AND session_generation=1 AND revoked_at IS NOT NULL",accountID).Scan(&count).Error)
 if count!=1{t.Fatal("old session not durably revoked",count)}
 me(fresh,200);me(unrelated,200)
 call("GET","/identity/public/auth/recovery/validate","",issue(entity.PasswordReset,2),nil,200)

 // Confirm actual signed relay retries after an acknowledgment failure;
 // receipt storage stays immutable and duplicate delivery is idempotent.
 waitFor(t,ctx,45*time.Second,func()bool{
  var n int64;err:=db.Raw("SELECT count(*) FROM audit.school_mutation_receipts WHERE school_id=?",school).Scan(&n).Error
  return err==nil&&n==2&&receiptAttempts.Load()>=3
 })
 raw:=jsonText(t,map[string]any{"schoolId":school,"operation":"DISABLE","pageNumber":1,"rowSize":25})
 req,err:=http.NewRequestWithContext(ctx,"POST",platform.URL+"/audit/restricted/school-mutation-receipts/search",bytes.NewBufferString(raw));must(t,err)
 token,err:=core.SignWorkload(readerKey.config.KeyID,readerKey.config.Issuer,readerKey.private,core.WorkloadClaims{Subject:readerKey.config.Subject,Audience:"platform",Scopes:readerKey.config.Scopes,ExpiresAt:time.Now().Add(time.Minute)});must(t,err)
 req.Header.Set("Authorization","Bearer "+token);req.Header.Set("Content-Type","application/json")
 resp,err:=client.Do(req);must(t,err);defer resp.Body.Close()
 var page struct{Records []struct{SchoolID string `json:"schoolId"`;Operation string `json:"operation"`;PreviousVersion string `json:"previousVersion"`} `json:"records"`}
 must(t,json.NewDecoder(resp.Body).Decode(&page))
 if resp.StatusCode!=200||len(page.Records)!=1||page.Records[0].SchoolID!=school||page.Records[0].Operation!="DISABLE"||page.Records[0].PreviousVersion!="1"{t.Fatal("incorrect filtered receipt",resp.StatusCode,page)}
 // Stop the real registry and verify the configured direct service boundary
 // fails closed. Restart with the same database and credentials.
 must(t,schoolProcess.Process.Kill());_ = schoolProcess.Wait()
 if _,err=integrationModule.Client().ListConnections(ctx,school,personID.String());err==nil{t.Fatal("registry outage accepted")}
 if _,err=lookup(ctx,school);err==nil{t.Fatal("registry outage returned a state")}
 _=startSchool()
 me(fresh,200);me(unrelated,200);mailbox(fresh,200)
}

package authentication

import (
 "context"
 corehttp "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http"
 "database/sql"
 "net/http"
 "net/http/httptest"
 "os"
 "strings"
 "testing"

 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/services/authentication"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/resources/migrations"
 coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
)

func TestOperatorAuthentication(t *testing.T){
 dsn:=os.Getenv("INTERNAL_TOOL_AUTH_TEST_DSN")
 if dsn==""{t.Skip("requires a fresh disposable database")}
 ctx:=context.Background()
 db,err:=coredatabase.OpenProjectPool(ctx,dsn);if err!=nil{t.Fatal(err)};defer db.Close()
 if err=migrations.Apply(ctx,db);err!=nil{t.Fatal(err)}
 if err=migrations.Apply(ctx,db);err!=nil{t.Fatal("migration not idempotent",err)}
 if err=migrations.Check(ctx,db);err!=nil{t.Fatal(err)}
 service:=authentication.New(db)
 if err=service.Bootstrap(ctx,"operator@example.com","test-only-password-123");err!=nil{t.Fatal(err)}
 if err=service.Bootstrap(ctx,"second@example.com","test-only-password-456");err==nil{t.Fatal("second bootstrap accepted")}
 controller,err:=New(service,"https://internal.example.com");if err!=nil{t.Fatal(err)}
 mux:=http.NewServeMux();controller.Register(mux)
 loginBody:=`{"email":"operator@example.com","password":"test-only-password-123"}`
 request:=func(method,path,body,origin string,cookies []*http.Cookie,csrf string)*httptest.ResponseRecorder{
  r:=httptest.NewRequest(method,path,strings.NewReader(body));if origin!=""{r.Header.Set("Origin",origin)};if csrf!=""{r.Header.Set("X-CSRF-Token",csrf)};for _,c:=range cookies{r.AddCookie(c)}
  w:=httptest.NewRecorder();corehttp.Ingress(mux).ServeHTTP(w,r);return w
 }
 const login="/internal/public/auth/login"
 const me="/internal/protected/auth/me"
 const logout="/internal/protected/auth/logout"
 if w:=request("POST",login,loginBody,"https://school.example.com",nil,"");w.Code!=403{t.Fatalf("wrong origin: %d",w.Code)}
 if w:=request("POST",login,`{"email":"operator@example.com","password":"wrong"}`,"https://internal.example.com",nil,"");w.Code!=401{t.Fatalf("wrong password: %d",w.Code)}
 w:=request("POST",login,loginBody,"https://internal.example.com",nil,"");if w.Code!=200{t.Fatalf("login failed: %d %s",w.Code,w.Body.String())}
 cookies:=w.Result().Cookies();if len(cookies)!=2{t.Fatal("expected session and CSRF cookies")}
 var csrf,token string
 for _,c:=range cookies{
  if !c.Secure||c.Path!="/"||c.Domain!=""||c.SameSite!=http.SameSiteStrictMode{t.Fatal("unsafe cookie")}
  switch c.Name{case sessionCookie:if !c.HttpOnly{t.Fatal("session not HttpOnly")};token=c.Value;case csrfCookie:csrf=c.Value}
 }
 if token==""||csrf==""||strings.Contains(w.Body.String(),token){t.Fatal("missing cookie or exposed token")}
 if w:=request("GET",me,"","",cookies,"");w.Code!=200||!strings.Contains(w.Body.String(),`"persona":"L2"`){t.Fatalf("me failed: %d",w.Code)}
 if w:=request("GET",me,"","",[]*http.Cookie{{Name:"__Host-gh_access",Value:token}},"");w.Code!=401{t.Fatal("School cookie accepted")}
 if w:=request("POST",logout,"","https://internal.example.com",cookies,"wrong");w.Code!=403{t.Fatal("wrong CSRF accepted")}
 if w:=request("POST",logout,"","https://internal.example.com",cookies,csrf);w.Code!=200{t.Fatalf("logout failed: %d",w.Code)}
 if w:=request("GET",me,"","",cookies,"");w.Code!=401{t.Fatal("revoked session accepted")}
 var count int
 if err:=db.QueryRowContext(ctx,"SELECT count(*) FROM internal_tool.workflow_events e JOIN internal_tool.audit_outbox o ON o.workflow_event_id=e.id WHERE e.action IN (\x27AUTH_LOGIN_SUCCEEDED\x27,\x27AUTH_LOGOUT_SUCCEEDED\x27) AND e.correlation_id IS NOT NULL").Scan(&count);err!=nil||count!=2{t.Fatal("authentication audit not committed",count,err)}
 if err:=db.QueryRowContext(ctx,"SELECT count(*) FROM internal_tool.sessions WHERE revoked_at IS NOT NULL").Scan(&count);err!=nil||count!=1{t.Fatal("revocation not persisted")}
 for i:=0;i<5;i++{if _,err:=service.Login(ctx,"unknown@example.com","wrong");err!=authentication.ErrInvalid{t.Fatalf("unknown account attempt %d: %v",i,err)}}
 if _,err:=service.Login(ctx,"unknown@example.com","wrong");err!=authentication.ErrLimited{t.Fatal("attempt budget not enforced")}
 if err:=db.QueryRowContext(ctx,"SELECT COALESCE(sum(attempts),0) FROM internal_tool.authentication_counters WHERE outcome=\x27REJECTED\x27").Scan(&count);err!=nil||count!=6{t.Fatal("rejected attempt counter",count,err)}
 // Operator administration uses a live L2 session and revokes disabled sessions.
 current,err:=service.Login(ctx,"operator@example.com","test-only-password-123");if err!=nil{t.Fatal(err)}
 created,err:=service.CreateOperator(ctx,current.Operator,"new-l1@example.com","initial-password-123","L1");if err!=nil{t.Fatal(err)}
 child,err:=service.Login(ctx,created.Email,"initial-password-123");if err!=nil{t.Fatal(err)}
 if _,err=service.CreateOperator(ctx,child.Operator,"unauthorized@example.com","initial-password-123","L2");err==nil{t.Fatal("L1 created operator")}
 if err=service.SetOperatorStatus(ctx,current.Operator,current.Operator.ID,"DISABLED");err==nil{t.Fatal("self-disable accepted")}
 if err=service.SetOperatorStatus(ctx,current.Operator,created.ID,"DISABLED");err!=nil{t.Fatal(err)}
 if _,err=service.Verify(ctx,child.Token);err!=authentication.ErrInvalid{t.Fatal("disabled operator session accepted")}
 if err=service.SetOperatorStatus(ctx,current.Operator,created.ID,"ACTIVE");err!=nil{t.Fatal(err)}
 if _,err=service.Verify(ctx,child.Token);err!=authentication.ErrInvalid{t.Fatal("old session restored on activation")}
 // Disable a live operator and ensure existing sessions are immediately invalid.
 result,err:=service.Login(ctx,"operator@example.com","test-only-password-123");if err!=nil{t.Fatal(err)}
 if _,err=db.ExecContext(ctx,"UPDATE internal_tool.operators SET status='DISABLED' WHERE email='operator@example.com'");err!=nil{t.Fatal(err)}
 if _,err=service.Verify(ctx,result.Token);err!=authentication.ErrInvalid{t.Fatal("disabled operator session accepted")}
}
func TestOriginConfiguration(t *testing.T){
 for _,origin:=range []string{"","http://example.com","https://example.com/path","https://user@example.com"}{if _,err:=New(authentication.New((*sql.DB)(nil)),origin);err==nil{t.Errorf("accepted invalid origin %q",origin)}}
}

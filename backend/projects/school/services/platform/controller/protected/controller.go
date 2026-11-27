package protected

import (
 "crypto/subtle"
 "encoding/json"
 "errors"
 audit "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/audit/client/activity"
 "io"
 "net/http"
 "strings"

 auth "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/identity/client/authentication"
 client "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/platform/client"
)

// Register exposes only explicitly enumerated operations. Platform remains
// responsible for live session verification, RBAC and resource authorization.
func Register(mux *http.ServeMux,platform *auth.Client,activity *audit.Client,origin string){
 for _,operation:=range client.Operations(){
  op:=operation
  mux.HandleFunc(op.Method+" "+op.Path,func(w http.ResponseWriter,r *http.Request){
   w.Header().Set("Cache-Control","no-store")
   target,ok:=op.Resolve(r)
   if !ok || r.URL.RawQuery!=""{http.Error(w,"invalid request",400);return}
   if r.Method!=http.MethodGet{
    csrf,err:=r.Cookie("gh_csrf")
    if origin==""||r.Header.Get("Origin")!=origin||err!=nil||csrf.Value==""||len(csrf.Value)>256||subtle.ConstantTimeCompare([]byte(csrf.Value),[]byte(r.Header.Get("X-CSRF-Token")))!=1{http.Error(w,"forbidden",403);return}
   }
   token,err:=r.Cookie("__Host-gh_access")
   if err!=nil||token.Value==""||len(token.Value)>2048||strings.ContainsAny(token.Value," \t\r\n,;"){http.Error(w,"unauthorized",401);return}
   var input json.RawMessage
   if r.Method!=http.MethodGet{
    body,err:=io.ReadAll(http.MaxBytesReader(w,r.Body,1<<20))
    if err!=nil{http.Error(w,"invalid request",400);return}
    if len(strings.TrimSpace(string(body)))>0{
     if !json.Valid(body)||!strings.HasPrefix(strings.TrimSpace(string(body)),"{"){http.Error(w,"invalid request",400);return}
     input=json.RawMessage(body)
    }
   }
   if op.Method=="POST" && op.Path=="/audit/protected/search" {
    result,err:=activity.Search(r.Context(),input,token.Value)
    if errors.Is(err,audit.ErrInvalidRequest){http.Error(w,"invalid audit request",400);return}
    if err!=nil{http.Error(w,"service unavailable",503);return}
    w.Header().Set("Content-Type","application/json")
    w.WriteHeader(result.Status)
    _,_=w.Write(result.Body)
    return
   }
   if strings.HasPrefix(op.Scope,"identity.") {
    id:=r.PathValue("accountId")
    if id=="" {id=r.PathValue("roleId")}
    if id=="" {id=r.PathValue("requestId")}
    key:=""
    switch op.Scope {
    case "identity.accounts.create","identity.approvals.lifecycle","identity.approvals.permissions","identity.approvals.assignments":
     values:=r.Header.Values("Idempotency-Key")
     if len(values)>1 {http.Error(w,"invalid idempotency key",400);return}
     if len(values)==1 {key=values[0]}
    }
    result,err:=platform.Business(r.Context(),op.Scope,id,input,token.Value,key)
    if errors.Is(err,auth.ErrInvalidBusinessRequest){http.Error(w,"invalid identity request",400);return}
    if err!=nil{http.Error(w,"service unavailable",503);return}
    w.Header().Set("Content-Type","application/json")
    w.WriteHeader(result.Status)
    _,_=w.Write(result.Body)
    return
   }
   result,err:=platform.CallBusiness(r.Context(),op.Method,target,input,token.Value)
   if err!=nil{http.Error(w,"service unavailable",503);return}
   w.Header().Set("Content-Type","application/json")
   w.WriteHeader(result.Status)
   _,_=w.Write(result.Body)
  })
 }
}

// RegisterOAuthCallback is the explicitly permitted provider ingress. The
// Platform consumes the one-time state and binds it to the pending connection.
func RegisterOAuthCallback(mux *http.ServeMux,platform *auth.Client){
 mux.HandleFunc("GET /integrations/public/connections/oauth/{provider}/callback",func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Cache-Control","no-store")
  w.Header().Set("Referrer-Policy","no-referrer")
  if len(r.URL.RawQuery)>8192{http.Error(w,"invalid OAuth callback",400);return}
  q:=r.URL.Query()
  provider:=strings.ToLower(r.PathValue("provider"))
  if (provider!="gmail"&&provider!="outlook")||q.Get("error")!=""||len(q["code"])!=1||len(q["state"])!=1||len(q.Get("code"))==0||len(q.Get("code"))>4096||len(q.Get("state"))!=43{http.Error(w,"invalid OAuth callback",400);return}
  for _,c:=range q.Get("state"){if !(c>='a'&&c<='z'||c>='A'&&c<='Z'||c>='0'&&c<='9'||c=='_'||c=='-'){http.Error(w,"invalid OAuth callback",400);return}}
  result,err:=platform.FinishOAuth(r.Context(),provider,q.Get("code"),q.Get("state"))
  if err!=nil{http.Error(w,"OAuth completion unavailable",503);return}
  w.Header().Set("Content-Type","application/json")
  w.WriteHeader(result.Status);_,_=w.Write(result.Body)
 })
}

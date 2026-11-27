package protected

import (
 "crypto/subtle"
 "encoding/json"
 "io"
 "net/http"
 "strings"

 auth "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/identity/client/authentication"
 client "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/platform/client"
)

// Register exposes only explicitly enumerated operations. Platform remains
// responsible for live session verification, RBAC and resource authorization.
func Register(mux *http.ServeMux,platform *auth.Client,origin string){
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
   result,err:=platform.CallBusiness(r.Context(),op.Method,target,input,token.Value)
   if err!=nil{http.Error(w,"service unavailable",503);return}
   w.Header().Set("Content-Type","application/json")
   w.WriteHeader(result.Status)
   _,_=w.Write(result.Body)
  })
 }
}

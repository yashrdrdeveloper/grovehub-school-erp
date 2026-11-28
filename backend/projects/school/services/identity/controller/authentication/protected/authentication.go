package protected

import (
 "crypto/subtle"
 "encoding/json"
 "io"
 "net/http"
 "strings"

 client "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/identity/client/authentication"
)

type Controller struct{ platform *client.Client; origin string }
func New(platform *client.Client,origin string)*Controller{return &Controller{platform:platform,origin:origin}}
func(c *Controller) Register(mux *http.ServeMux){
 mux.HandleFunc("POST /identity/protected/auth/logout",c.Logout)
 mux.HandleFunc("POST /identity/protected/auth/logout-all",c.LogoutAll)
 mux.HandleFunc("POST /identity/protected/auth/change-password",c.ChangePassword)
 mux.HandleFunc("POST /identity/protected/auth/reauthenticate",c.Reauthenticate)
 mux.HandleFunc("GET /identity/protected/auth/sessions",c.Sessions)
 mux.HandleFunc("DELETE /identity/protected/auth/sessions/{sessionId}",c.RevokeSession)
 mux.HandleFunc("GET /identity/protected/auth/me",c.Me)
}
func(c *Controller) proof(w http.ResponseWriter,r *http.Request)(string,bool){
 w.Header().Set("Cache-Control","no-store")
 if r.Method!=http.MethodGet {
  csrf,err:=r.Cookie("gh_csrf")
  if c.origin=="" || r.Header.Get("Origin")!=c.origin || err!=nil || csrf.Value=="" || len(csrf.Value)>256 || subtle.ConstantTimeCompare([]byte(csrf.Value),[]byte(r.Header.Get("X-CSRF-Token")))!=1 {http.Error(w,"forbidden",403);return "",false}
 }
 token,err:=r.Cookie("__Host-gh_access")
 if err!=nil || token.Value=="" || len(token.Value)>2048 || strings.ContainsAny(token.Value," \t\r\n,;"){http.Error(w,"unauthorized",401);return "",false}
 return token.Value,true
}
func decode(w http.ResponseWriter,r *http.Request,out any)bool{
 d:=json.NewDecoder(http.MaxBytesReader(w,r.Body,16<<10));d.DisallowUnknownFields()
 if err:=d.Decode(out);err!=nil{http.Error(w,"invalid request",400);return false}
 var extra any
 if err:=d.Decode(&extra);err!=io.EOF{http.Error(w,"invalid request",400);return false}
 return true
}
func write(w http.ResponseWriter,result client.Result,err error){
 if err!=nil{http.Error(w,"authentication unavailable",503);return}
 // Protected endpoints may only clear existing session cookies.
 if result.Status>=200&&result.Status<300{
  for _,cookie:=range result.Cookies {
   if cookie.MaxAge>=0{continue}
   path:="/";httpOnly:=true
   switch cookie.Name{
   case "__Host-gh_access":
   case "__Secure-gh_refresh":path="/identity/public/auth/refresh"
   case "gh_csrf":httpOnly=false
   default:continue
   }
   http.SetCookie(w,&http.Cookie{Name:cookie.Name,Path:path,MaxAge:-1,Secure:true,HttpOnly:httpOnly,SameSite:http.SameSiteStrictMode})
  }
 }
 w.Header().Set("Content-Type","application/json");w.WriteHeader(result.Status);_,_=w.Write(result.Body)
}
func(c *Controller) Logout(w http.ResponseWriter,r *http.Request){
 proof,ok:=c.proof(w,r);if !ok{return}
 
 result,err:=c.platform.Logout(r.Context(),proof)
 write(w,result,err)
}
func(c *Controller) LogoutAll(w http.ResponseWriter,r *http.Request){
 proof,ok:=c.proof(w,r);if !ok{return}
 
 result,err:=c.platform.LogoutAll(r.Context(),proof)
 write(w,result,err)
}
func(c *Controller) ChangePassword(w http.ResponseWriter,r *http.Request){
 proof,ok:=c.proof(w,r);if !ok{return}
 var input client.ChangePassword;if !decode(w,r,&input){return}
 result,err:=c.platform.ChangePassword(r.Context(),input,proof)
 write(w,result,err)
}
func(c *Controller) Reauthenticate(w http.ResponseWriter,r *http.Request){
 proof,ok:=c.proof(w,r);if !ok{return}
 var input client.Reauthenticate;if !decode(w,r,&input){return}
 result,err:=c.platform.Reauthenticate(r.Context(),input,proof)
 write(w,result,err)
}
func(c *Controller) Sessions(w http.ResponseWriter,r *http.Request){
 proof,ok:=c.proof(w,r);if !ok{return}
 
 result,err:=c.platform.Sessions(r.Context(),proof)
 write(w,result,err)
}
func(c *Controller) RevokeSession(w http.ResponseWriter,r *http.Request){
 proof,ok:=c.proof(w,r);if !ok{return}
 id:=r.PathValue("sessionId");if !client.ValidSessionID(id){http.Error(w,"invalid session ID",400);return}
 result,err:=c.platform.RevokeSession(r.Context(),id,proof)
 write(w,result,err)
}
func(c *Controller) Me(w http.ResponseWriter,r *http.Request){
 proof,ok:=c.proof(w,r);if !ok{return}
 
 result,err:=c.platform.Me(r.Context(),proof)
 write(w,result,err)
}

package public

import (
 "encoding/json"
 "io"
 "net/http"
 "sync"
 "strings"
 "crypto/subtle"
 "time"

 client "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/identity/client/authentication"
)

type Controller struct {
 platform *client.Client
 origin string
 mu sync.Mutex
 minute time.Time
 counts map[string]int
}

func New(platform *client.Client,origin string)*Controller {
 return &Controller{platform:platform,origin:origin,counts:make(map[string]int)}
}
func (c *Controller) Register(mux *http.ServeMux) {
 mux.HandleFunc("POST /identity/public/auth/login",c.Login)
 mux.HandleFunc("POST /identity/public/auth/forgot-password",c.ForgotPassword)
 mux.HandleFunc("POST /identity/public/auth/refresh",c.Refresh)
 mux.HandleFunc("GET /identity/public/auth/recovery/validate",c.RecoveryValidate)
 mux.HandleFunc("GET /identity/public/auth/invitations/validate",c.InvitationValidate)
 mux.HandleFunc("POST /identity/public/auth/recovery/exchange",c.RecoveryExchange)
 mux.HandleFunc("POST /identity/public/auth/invitations/exchange",c.InvitationExchange)
 mux.HandleFunc("GET /identity/public/auth/reset-password/status",c.ResetStatus)
 mux.HandleFunc("POST /identity/public/auth/reset-password",c.ResetPassword)
 mux.HandleFunc("GET /identity/public/accounts/setup/status",c.SetupStatus)
 mux.HandleFunc("POST /identity/public/accounts/setup",c.Setup)
}
// A deliberately bounded global migration limiter; replace with a shared,
// trusted-client-aware limiter before deploying multiple instances.
func (c *Controller) allow(w http.ResponseWriter,r *http.Request,operation string,limit int)bool {
 w.Header().Set("Cache-Control","no-store")
 if c.origin=="" || (r.Method!=http.MethodGet && r.Header.Get("Origin")!=c.origin) {
  http.Error(w,"forbidden",http.StatusForbidden);return false
 }
 c.mu.Lock()
 defer c.mu.Unlock()
 now:=time.Now().Truncate(time.Minute)
 if !now.Equal(c.minute) {c.minute=now;c.counts=make(map[string]int)}
 if c.counts[operation]>=limit {w.Header().Set("Retry-After","60");http.Error(w,"too many requests",http.StatusTooManyRequests);return false}
 c.counts[operation]++
 return true
}
func decode(w http.ResponseWriter,r *http.Request,out any)bool {
 d:=json.NewDecoder(http.MaxBytesReader(w,r.Body,16<<10))
 d.DisallowUnknownFields()
 if err:=d.Decode(out);err!=nil {http.Error(w,"invalid request",400);return false}
 var extra any
 if err:=d.Decode(&extra);err!=io.EOF {http.Error(w,"invalid request",400);return false}
 return true
}
func (c *Controller) Login(w http.ResponseWriter,r *http.Request) {
 if !c.allow(w,r,"login",12) {return}
 var input client.Login
 if !decode(w,r,&input) {return}
 result,err:=c.platform.Login(r.Context(),input)
 if err!=nil {http.Error(w,"authentication unavailable",503);return}
 if result.Status==http.StatusOK {
  // Reissue only the known School session cookies with fixed attributes.
  for _, cookie:=range result.Cookies {
   switch cookie.Name {
   case "__Host-gh_access":
    http.SetCookie(w,&http.Cookie{Name:cookie.Name,Value:cookie.Value,Path:"/",MaxAge:900,Secure:true,HttpOnly:true,SameSite:http.SameSiteLaxMode})
   case "__Secure-gh_refresh":
    http.SetCookie(w,&http.Cookie{Name:cookie.Name,Value:cookie.Value,Path:"/identity/public/auth/refresh",MaxAge:2592000,Secure:true,HttpOnly:true,SameSite:http.SameSiteStrictMode})
   case "gh_csrf":
    http.SetCookie(w,&http.Cookie{Name:cookie.Name,Value:cookie.Value,Path:"/",MaxAge:2592000,Secure:true,SameSite:http.SameSiteStrictMode})
   }
  }
 }
 write(w,result)
}
func (c *Controller) ForgotPassword(w http.ResponseWriter,r *http.Request) {
 if !c.allow(w,r,"forgot-password",5) {return}
 var input client.ForgotPassword
 if !decode(w,r,&input) {return}
 result,err:=c.platform.ForgotPassword(r.Context(),input)
 if err!=nil {http.Error(w,"recovery unavailable",503);return}
 write(w,result)
}
func write(w http.ResponseWriter,result client.Result) {
 w.Header().Set("Content-Type","application/json")
 w.WriteHeader(result.Status)
 _,_=w.Write(result.Body)
}

func(c *Controller) proof(w http.ResponseWriter,r *http.Request,operation,cookie,prefix string,csrf bool,limit int)(string,bool){
 if !c.allow(w,r,operation,limit){return "",false}
 if csrf {
  token,err:=r.Cookie("gh_csrf")
  if err!=nil || len(token.Value)==0 || len(token.Value)>256 || subtle.ConstantTimeCompare([]byte(token.Value),[]byte(r.Header.Get("X-CSRF-Token")))!=1 {http.Error(w,"forbidden",403);return "",false}
 }
 var proof string
 if cookie!="" {if v,err:=r.Cookie(cookie);err==nil {proof=v.Value}} else {
  values:=r.Header.Values("Authorization")
  if len(values)==1 {proof,_=strings.CutPrefix(values[0],prefix+" ")}
 }
 if proof=="" || len(proof)>2048 || strings.ContainsAny(proof," \t\r\n,;") {http.Error(w,"invalid credential",401);return "",false}
 return proof,true
}
func writeProofResult(w http.ResponseWriter,result client.Result,err error){
 if err!=nil {http.Error(w,"authentication unavailable",503);return}
 for _, cookie:=range result.Cookies {
  path:="";httpOnly:=true;maxAge:=cookie.MaxAge;sameSite:=http.SameSiteStrictMode
  switch cookie.Name {
  case "__Host-gh_access":path="/";sameSite=http.SameSiteLaxMode
  case "__Secure-gh_refresh":path="/identity/public/auth/refresh"
  case "__Secure-gh_reset":path="/identity/public/auth/reset-password"
  case "__Secure-gh_setup":path="/identity/public/accounts/setup"
  case "gh_csrf":path="/";httpOnly=false
  default:continue
  }
  if maxAge==0 {continue}
  http.SetCookie(w,&http.Cookie{Name:cookie.Name,Value:cookie.Value,Path:path,MaxAge:maxAge,Secure:true,HttpOnly:httpOnly,SameSite:sameSite})
 }
 write(w,result)
}
func(c *Controller) Refresh(w http.ResponseWriter,r *http.Request){
 proof,ok:=c.proof(w,r,"refresh","__Secure-gh_refresh","",true,30)
 if !ok{return}
 
 result,err:=c.platform.Refresh(r.Context(),proof)
 writeProofResult(w,result,err)
}
func(c *Controller) RecoveryValidate(w http.ResponseWriter,r *http.Request){
 proof,ok:=c.proof(w,r,"recovery.validate","","Recovery",false,20)
 if !ok{return}
 
 result,err:=c.platform.RecoveryValidate(r.Context(),proof)
 writeProofResult(w,result,err)
}
func(c *Controller) InvitationValidate(w http.ResponseWriter,r *http.Request){
 proof,ok:=c.proof(w,r,"invitation.validate","","Invitation",false,20)
 if !ok{return}
 
 result,err:=c.platform.InvitationValidate(r.Context(),proof)
 writeProofResult(w,result,err)
}
func(c *Controller) RecoveryExchange(w http.ResponseWriter,r *http.Request){
 proof,ok:=c.proof(w,r,"recovery.exchange","","Recovery",false,8)
 if !ok{return}
 
 result,err:=c.platform.RecoveryExchange(r.Context(),proof)
 writeProofResult(w,result,err)
}
func(c *Controller) InvitationExchange(w http.ResponseWriter,r *http.Request){
 proof,ok:=c.proof(w,r,"invitation.exchange","","Invitation",false,8)
 if !ok{return}
 
 result,err:=c.platform.InvitationExchange(r.Context(),proof)
 writeProofResult(w,result,err)
}
func(c *Controller) ResetStatus(w http.ResponseWriter,r *http.Request){
 proof,ok:=c.proof(w,r,"reset.status","__Secure-gh_reset","",false,20)
 if !ok{return}
 
 result,err:=c.platform.ResetStatus(r.Context(),proof)
 writeProofResult(w,result,err)
}
func(c *Controller) ResetPassword(w http.ResponseWriter,r *http.Request){
 proof,ok:=c.proof(w,r,"reset","__Secure-gh_reset","",true,20)
 if !ok{return}
 var input client.ResetPassword;if !decode(w,r,&input){return}
 result,err:=c.platform.ResetPassword(r.Context(),input,proof)
 writeProofResult(w,result,err)
}
func(c *Controller) SetupStatus(w http.ResponseWriter,r *http.Request){
 proof,ok:=c.proof(w,r,"setup.status","__Secure-gh_setup","",false,20)
 if !ok{return}
 
 result,err:=c.platform.SetupStatus(r.Context(),proof)
 writeProofResult(w,result,err)
}
func(c *Controller) Setup(w http.ResponseWriter,r *http.Request){
 proof,ok:=c.proof(w,r,"setup","__Secure-gh_setup","",true,20)
 if !ok{return}
 var input client.Setup;if !decode(w,r,&input){return}
 result,err:=c.platform.Setup(r.Context(),input,proof)
 writeProofResult(w,result,err)
}

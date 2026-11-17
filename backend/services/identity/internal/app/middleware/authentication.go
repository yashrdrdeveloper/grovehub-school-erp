package middleware
import(
 "crypto/subtle"
 "net/http"
 "strings"
 "time"
 httpresponse "grovehub/libs/core/http/response"
 coreerrors "grovehub/services/identity/internal/app/constants/error"
 accountrepo "grovehub/services/identity/internal/app/repository/account"
 sessionrepo "grovehub/services/identity/internal/app/repository/session"
 "grovehub/services/identity/internal/app/security/principal"
 helpers "grovehub/services/identity/internal/app/services/authentication/helpers"
 "grovehub/services/identity/internal/app/entity"
)
type Authentication struct{Accounts accountrepo.Repository;Sessions sessionrepo.Repository;AdminIDs map[string]bool;AllowedOrigin string}
const AccessCookie="__Host-gh_access"
const RefreshCookie="__Secure-gh_refresh"
const CSRFName="gh_csrf"
const ResetCookie="__Secure-gh_reset"
const SetupCookie="__Secure-gh_setup"
func (m Authentication) Forbidden(w http.ResponseWriter,r *http.Request){_ = httpresponse.WriteErrorContext(r.Context(),w,403,coreerrors.Forbidden)}
func (m Authentication) Unauthorized(w http.ResponseWriter,r *http.Request){_ = httpresponse.WriteErrorContext(r.Context(),w,401,coreerrors.Unauthorized)}
func (m Authentication) Origin(next http.Handler)http.Handler{
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.Method!="GET"&&r.Method!="HEAD"&&r.Method!="OPTIONS"{
   if m.AllowedOrigin==""||r.Header.Get("Origin")!=m.AllowedOrigin{m.Forbidden(w,r);return}
  }
  next.ServeHTTP(w,r)
 })
}
func (m Authentication) CSRF(next http.Handler)http.Handler{
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.Method!="GET"&&r.Method!="HEAD"&&r.Method!="OPTIONS"{
   c,err:=r.Cookie(CSRFName)
   header:=r.Header.Get("X-CSRF-Token")
   if err!=nil||len(header)==0||len(header)>256||subtle.ConstantTimeCompare([]byte(c.Value),[]byte(header))!=1{m.Forbidden(w,r);return}
  }
  next.ServeHTTP(w,r)
 })
}
func (m Authentication) Protected(next http.Handler)http.Handler{
 return m.Origin(m.CSRF(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  c,err:=r.Cookie(AccessCookie);if err!=nil{m.Unauthorized(w,r);return}
  hash,err:=helpers.Digest(c.Value,32);if err!=nil{m.Unauthorized(w,r);return}
  s,err:=m.Sessions.FindAccess(r.Context(),hash,time.Now().UTC());if err!=nil{m.Unauthorized(w,r);return}
  account,err:=m.Accounts.FindByID(r.Context(),s.AccountID);if err!=nil||account.Status!=entity.Active{m.Unauthorized(w,r);return}
  tenant:="";if account.TenantID!=nil{tenant=*account.TenantID}
  authenticated:=s.CreatedAt;if s.ReauthenticatedAt!=nil{authenticated=*s.ReauthenticatedAt}
  p:=principal.Principal{AccountID:account.AccountID,PersonGID:account.PersonGID,TenantID:tenant,SessionID:s.SessionID,Admin:m.AdminIDs[account.AccountID],AuthenticatedAt:authenticated}
  next.ServeHTTP(w,r.WithContext(principal.With(r.Context(),p)))
 })))
}
func IssueCookies(w http.ResponseWriter,access,refresh,csrf string){
 http.SetCookie(w,&http.Cookie{Name:AccessCookie,Value:access,Path:"/",MaxAge:900,Secure:true,HttpOnly:true,SameSite:http.SameSiteLaxMode})
 http.SetCookie(w,&http.Cookie{Name:RefreshCookie,Value:refresh,Path:"/identity/public/auth/refresh",MaxAge:2592000,Secure:true,HttpOnly:true,SameSite:http.SameSiteStrictMode})
 if csrf!=""{http.SetCookie(w,&http.Cookie{Name:CSRFName,Value:csrf,Path:"/",MaxAge:2592000,Secure:true,HttpOnly:false,SameSite:http.SameSiteStrictMode})}
 w.Header().Set("Cache-Control","no-store")
}
func ClearSessionCookies(w http.ResponseWriter){
 for _,c:=range []struct{name,path string}{{AccessCookie,"/"},{RefreshCookie,"/identity/public/auth/refresh"},{CSRFName,"/"}}{
  http.SetCookie(w,&http.Cookie{Name:c.name,Path:c.path,MaxAge:-1,Secure:true,HttpOnly:c.name!=CSRFName,SameSite:http.SameSiteStrictMode})
 }
 w.Header().Set("Cache-Control","no-store")
}
func ReadCookie(r *http.Request,name string)string{c,err:=r.Cookie(name);if err!=nil{return ""};return c.Value}
func SetGrant(w http.ResponseWriter,name,path,value string){http.SetCookie(w,&http.Cookie{Name:name,Path:path,Value:value,MaxAge:600,Secure:true,HttpOnly:true,SameSite:http.SameSiteStrictMode});w.Header().Set("Cache-Control","no-store")}
func ClearGrant(w http.ResponseWriter,name,path string){SetGrant(w,name,path,"");http.SetCookie(w,&http.Cookie{Name:name,Path:path,MaxAge:-1,Secure:true,HttpOnly:true,SameSite:http.SameSiteStrictMode})}
func Proof(r *http.Request,prefix string)string{
 auth:=r.Header.Get("Authorization")
 if !strings.HasPrefix(auth,prefix+" "){return ""}
 return strings.TrimPrefix(auth,prefix+" ")
}

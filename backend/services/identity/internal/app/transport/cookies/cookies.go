package cookies

import (
	"net/http"
	"strings"

	coregateway "grovehub/libs/core/gateway"
)

const (
	AccessCookie = coregateway.AccessCookie
	RefreshCookie = coregateway.RefreshCookie
	CSRFName = coregateway.CSRFName
	ResetCookie = coregateway.ResetCookie
	SetupCookie = coregateway.SetupCookie
)

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

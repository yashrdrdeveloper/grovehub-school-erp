package public

import (
 "encoding/json"
 "io"
 "net/http"
 "sync"
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
}
// A deliberately bounded global migration limiter; replace with a shared,
// trusted-client-aware limiter before deploying multiple instances.
func (c *Controller) allow(w http.ResponseWriter,r *http.Request,operation string,limit int)bool {
 w.Header().Set("Cache-Control","no-store")
 if c.origin=="" || r.Header.Get("Origin")!=c.origin {
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

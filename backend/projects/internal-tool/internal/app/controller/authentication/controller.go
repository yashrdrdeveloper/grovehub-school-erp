package authentication

import (
 "encoding/json"
 "errors"
 "io"
 "net/http"
 "net/url"
 "strings"

 service "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/services/authentication"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)

const sessionCookie="__Host-gh_internal_session"
const csrfCookie="__Host-gh_internal_csrf"

type Controller struct{service *service.Service;origin string}
func New(s *service.Service,origin string)(*Controller,error){
 u,err:=url.Parse(origin)
 if err!=nil||u.Host==""||u.User!=nil||u.Path!=""||u.RawQuery!=""||u.Fragment!=""||(u.Scheme!="https"&&!(u.Scheme=="http"&&(u.Hostname()=="127.0.0.1"||u.Hostname()=="localhost"))){return nil,errors.New("invalid INTERNAL_TOOL_ALLOWED_ORIGIN")}
 return &Controller{s,origin},nil
}
func(c *Controller) Register(mux *http.ServeMux){
 mux.HandleFunc("POST /internal/public/auth/login",c.login)
 mux.HandleFunc("GET /internal/protected/auth/me",c.me)
 mux.HandleFunc("POST /internal/protected/auth/logout",c.logout)
}
func write(w http.ResponseWriter,status int,data any){
 w.Header().Set("Content-Type","application/json")
 w.Header().Set("Cache-Control","no-store")
 w.WriteHeader(status);_ = json.NewEncoder(w).Encode(data)
}
func fail(w http.ResponseWriter,err error){
 status:=http.StatusServiceUnavailable;code:="UNAVAILABLE"
 if errors.Is(err,service.ErrInvalid){status=401;code="UNAUTHORIZED"}
 if errors.Is(err,service.ErrLimited){status=429;code="RATE_LIMITED";w.Header().Set("Retry-After","300")}
 write(w,status,map[string]any{"success":false,"errorCode":code})
}
func(c *Controller) mutationAllowed(r *http.Request)bool{return r.Header.Get("Origin")==c.origin&&r.URL.RawQuery==""}
func cookie(r *http.Request,name string)string{value:="";count:=0;for _,c:=range r.Cookies(){if c.Name==name{value=c.Value;count++}};if count!=1{return ""};return value}
func setCookie(w http.ResponseWriter,name,value string,httpOnly bool,maxAge int){
 http.SetCookie(w,&http.Cookie{Name:name,Value:value,Path:"/",Secure:true,HttpOnly:httpOnly,SameSite:http.SameSiteStrictMode,MaxAge:maxAge})
}
func(c *Controller) login(w http.ResponseWriter,r *http.Request){
 if !c.mutationAllowed(r){write(w,403,map[string]any{"success":false,"errorCode":"FORBIDDEN"});return}
 var input struct{Email string `json:"email"`;Password string `json:"password"`}
 decoder:=json.NewDecoder(http.MaxBytesReader(w,r.Body,4096));decoder.DisallowUnknownFields()
 var extra any
 if decoder.Decode(&input)!=nil||decoder.Decode(&extra)!=io.EOF{write(w,400,map[string]any{"success":false,"errorCode":"INVALID_REQUEST"});return}
 result,err:=c.service.Login(r.Context(),input.Email,input.Password);if err!=nil{fail(w,err);return}
 ttl:=int(service.SessionTTL.Seconds())
 setCookie(w,sessionCookie,result.Token,true,ttl);setCookie(w,csrfCookie,result.CSRF,false,ttl)
 write(w,200,map[string]any{"success":true,"data":result.Operator})
}
func(c *Controller) me(w http.ResponseWriter,r *http.Request){
 if r.URL.RawQuery!=""{write(w,400,map[string]any{"success":false,"errorCode":"INVALID_REQUEST"});return}
 op,err:=c.service.Verify(r.Context(),cookie(r,sessionCookie));if err!=nil{fail(w,err);return}
 write(w,200,map[string]any{"success":true,"data":op})
}
func(c *Controller) logout(w http.ResponseWriter,r *http.Request){
 values:=r.Header.Values("X-CSRF-Token")
 csrf:=cookie(r,csrfCookie)
 if !c.mutationAllowed(r)||len(values)!=1||csrf==""||strings.TrimSpace(values[0])!=csrf{write(w,403,map[string]any{"success":false,"errorCode":"FORBIDDEN"});return}
 if err:=c.service.Logout(r.Context(),cookie(r,sessionCookie),csrf);err!=nil{fail(w,err);return}
 setCookie(w,sessionCookie,"",true,-1);setCookie(w,csrfCookie,"",false,-1)
 write(w,200,map[string]any{"success":true,"data":map[string]string{"status":"LOGGED_OUT"}})
}

// Protected resolves the live operator from the dedicated session cookie.
// Mutations require both the origin and session-bound CSRF proof.
func(c *Controller) Protected(mutate bool,next func(http.ResponseWriter,*http.Request,entity.Operator))http.HandlerFunc{
 return func(w http.ResponseWriter,r *http.Request){
  if r.URL.RawQuery!=""{write(w,400,map[string]any{"success":false,"errorCode":"INVALID_REQUEST"});return}
  var op entity.Operator;var err error
  if mutate{
   values:=r.Header.Values("X-CSRF-Token");csrf:=cookie(r,csrfCookie)
   if !c.mutationAllowed(r)||len(values)!=1||csrf==""||values[0]!=csrf{write(w,403,map[string]any{"success":false,"errorCode":"FORBIDDEN"});return}
   op,err=c.service.VerifyMutation(r.Context(),cookie(r,sessionCookie),csrf)
  }else{op,err=c.service.Verify(r.Context(),cookie(r,sessionCookie))}
  if err!=nil{fail(w,err);return}
  next(w,r,op)
 }
}

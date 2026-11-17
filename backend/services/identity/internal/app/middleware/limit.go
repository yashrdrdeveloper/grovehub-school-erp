package middleware

import(
 "net"
 "net/http"
 "sync"
 "time"
 httpresponse "grovehub/libs/core/http/response"
 definitions "grovehub/services/identity/internal/app/constants/error"
)

// Limiter is a bounded per-instance defense at the HTTP edge. Production
// multi-replica deployments must additionally enforce distributed limits.
type Limiter struct{
 mu sync.Mutex
 attempts map[string]attempt
}
type attempt struct{Count int;WindowEnd time.Time}
func NewLimiter()*Limiter{return &Limiter{attempts:make(map[string]attempt)}}
func(l *Limiter)Allowed(key string,max int,window time.Duration)bool{
 now:=time.Now()
 l.mu.Lock();defer l.mu.Unlock()
 if len(l.attempts)>100000{
  for k,v:=range l.attempts{if now.After(v.WindowEnd){delete(l.attempts,k)}}
 }
 if len(l.attempts)>100000{return false}
 v:=l.attempts[key]
 if now.After(v.WindowEnd){v=attempt{WindowEnd:now.Add(window)}}
 if v.Count>=max{return false}
 v.Count++;l.attempts[key]=v;return true
}
func(l *Limiter)Guard(route string,max int,window time.Duration,next http.Handler)http.Handler{
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  address,_,err:=net.SplitHostPort(r.RemoteAddr)
  if err!=nil{address=r.RemoteAddr}
  // Never trust X-Forwarded-For unless a trusted reverse-proxy middleware
  // has already established and authenticated the source address.
  key:=route+"|"+address
  if !l.Allowed(key,max,window){w.Header().Set("Retry-After","60");_ = httpresponse.WriteErrorContext(r.Context(),w,429,definitions.RateLimited);return}
  next.ServeHTTP(w,r)
 })
}

package security_test
import(
 "net/http"
 "net/http/httptest"
 "testing"
 middleware "grovehub/services/identity/internal/app/middleware"
)
func TestCrossOriginWriteRejected(t *testing.T){
 m:=middleware.Authentication{AllowedOrigin:"https://app.example.test"}
 called:=false
 h:=m.Origin(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){called=true;w.WriteHeader(204)}))
 req:=httptest.NewRequest(http.MethodPost,"/identity/public/auth/login",nil)
 req.Header.Set("Origin","https://evil.example.test")
 w:=httptest.NewRecorder();h.ServeHTTP(w,req)
 if called||w.Code!=403{t.Fatalf("cross origin write permitted: %d",w.Code)}
}

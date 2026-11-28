package controller_test
import (
 "net/http"
 "net/http/httptest"
 . "github.com/onsi/ginkgo/v2"
 . "github.com/onsi/gomega"
 "grovehub/launcher/internal/src/controller/impl"
 "grovehub/launcher/internal/src/routes"
)
// @Testing suite=LauncherHealth
func HealthSpecs(){
 Describe("HealthController",func(){
  It("registers generated probes and preserves readiness semantics",func(){
   ready:=false
   mux:=http.NewServeMux()
   routes.RegisterHTTP(mux,impl.New(func()bool{return ready}))
   run:=func(path string)*httptest.ResponseRecorder{w:=httptest.NewRecorder();mux.ServeHTTP(w,httptest.NewRequest("GET",path,nil));return w}
   Expect(run("/health").Code).To(Equal(200))
   Expect(run("/ready").Code).To(Equal(503))
   ready=true
   response:=run("/ready")
   Expect(response.Code).To(Equal(200))
   Expect(response.Header().Get("Cache-Control")).To(Equal("no-store"))
  })
 })
}

package impl_test

import (
 "net/http"
 "net/http/httptest"
 . "github.com/onsi/ginkgo/v2"
 . "github.com/onsi/gomega"
 "grovehub/launcher/internal/src/services/gateway/operations/ratelimit/impl"
)

// @Testing suite=LauncherRateLimit
func RateLimitSpecs(){
 Describe("Allow",func(){
  It("bounds each configured route and caller independently",func(){
   limiter:=impl.New(map[string]int{"POST /login":2,"GET /disabled":0})
   request:=httptest.NewRequest(http.MethodPost,"/login",nil)
   request.RemoteAddr="127.0.0.1:1000"
   Expect(limiter.Allow("/login",request)).To(BeTrue())
   request.RemoteAddr="127.0.0.1:2000"
   Expect(limiter.Allow("/login",request)).To(BeTrue())
   Expect(limiter.Allow("/login",request)).To(BeFalse())
   request.RemoteAddr="127.0.0.2:2000"
   Expect(limiter.Allow("/login",request)).To(BeTrue())
   Expect(limiter.Allow("/unknown",request)).To(BeTrue())
  })
 })
}

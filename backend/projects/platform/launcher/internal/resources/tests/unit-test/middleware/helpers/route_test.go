package helpers_test

import (
 "net/http"
 "net/http/httptest"
 . "github.com/onsi/ginkgo/v2"
 . "github.com/onsi/gomega"
 "grovehub/launcher/internal/src/middleware/helpers"
)

// @Testing suite=LauncherResponseRecorder
func RouteSpecs(){
 Describe("Route",func(){
  It("normalizes registered and unmatched paths",func(){
   r:=httptest.NewRequest(http.MethodGet,"/unknown",nil)
   Expect(helpers.Route(r,200)).To(Equal("unmatched"))
   r.Pattern="GET /"
   Expect(helpers.Route(r,404)).To(Equal("unmatched"))
   Expect(helpers.Route(r,200)).To(Equal("/"))
   r.Pattern="GET /accounts/{id}"
   Expect(helpers.Route(r,200)).To(Equal("/accounts/{id}"))
  })
 })
}

package impl_test

import (
 "context"
 "net/http"
 "net/http/httptest"

 . "github.com/onsi/ginkgo/v2"
 . "github.com/onsi/gomega"
 "grovehub/launcher/internal/src/services/gateway/impl"
 "grovehub/launcher/internal/src/services/gateway/operations/inprocess"
)

// @Testing suite=LauncherGatewayComposition
func GatewaySpecs(){
 Describe("Gateway composition",func(){
  It("exposes independent external and in-process boundaries and rejects missing internal authorization",func(){
   mux:=http.NewServeMux()
   mux.HandleFunc("GET /health",func(w http.ResponseWriter,_ *http.Request){w.WriteHeader(204)})
   gateway:=impl.New(mux,nil,nil,"",nil)
   response:=httptest.NewRecorder()
   gateway.External().ServeHTTP(response,httptest.NewRequest("GET","/health",nil))
   Expect(response.Code).To(Equal(204))
   Expect(gateway.Internal()).NotTo(BeNil())
   result,err:=inprocess.Invoke[string](context.Background(),nil,"read",func(context.Context)(string,error){Fail("unconfigured gateway invoked");return "",nil})
   Expect(result).To(BeEmpty())
   Expect(err).To(MatchError("internal gateway is not configured"))
  })
 })
}

package controller_test

import (
 "context"
 "net/http"
 "net/http/httptest"
 "strings"

 . "github.com/onsi/ginkgo/v2"
 . "github.com/onsi/gomega"

 integrationcontroller "grovehub/services/integration/controller/impl"
 integrationrequest "grovehub/services/integration/dto/request"
 integrationresponse "grovehub/services/integration/dto/response"
 "grovehub/services/integration/providers"
)

type serviceStub struct{}
func (*serviceStub) Providers(context.Context) []providers.Descriptor {
 return []providers.Descriptor{{ID:"GMAIL"}}
}
func (*serviceStub) Create(context.Context,string,integrationrequest.CreateConnection)(integrationresponse.Connection,error){
 return integrationresponse.Connection{},nil
}
func (*serviceStub) Get(context.Context,string,string)(integrationresponse.Connection,error){
 return integrationresponse.Connection{},nil
}
func (*serviceStub) List(context.Context,string)([]integrationresponse.Connection,error){
 return nil,nil
}
func (*serviceStub) BeginOAuth(context.Context,string,string)(integrationresponse.OAuthStart,error){
 return integrationresponse.OAuthStart{},nil
}
func (*serviceStub) FinishOAuth(context.Context,string,string)(integrationresponse.OAuthCompletion,error){
 return integrationresponse.OAuthCompletion{},nil
}
func (*serviceStub) RunTest(context.Context,string,string,string)(integrationresponse.ConnectionTest,error){
 return integrationresponse.ConnectionTest{},nil
}
func (*serviceStub) ListTests(context.Context,string,string)([]integrationresponse.ConnectionTest,error){
 return nil,nil
}
func (*serviceStub) Disconnect(context.Context,string,string)(integrationresponse.Connection,error){
 return integrationresponse.Connection{},nil
}

var _ = Describe("Integration HTTP authorization",func(){
 It("allows public provider discovery",func(){
  c:=integrationcontroller.NewConnectionController(&serviceStub{})
  recorder:=httptest.NewRecorder()
  req:=httptest.NewRequest(http.MethodGet,"/integrations/providers",nil)
  c.Providers(recorder,req)
  Expect(recorder.Code).To(Equal(http.StatusOK))
  Expect(recorder.Body.String()).To(ContainSubstring("GMAIL"))
 })
 It("denies unauthorized OAuth starts and mail tests",func(){
  c:=integrationcontroller.NewConnectionController(&serviceStub{})
  for _,testCase:=range []struct{route string;handler http.HandlerFunc}{
   {"/integrations/connections/abc/oauth/start",c.BeginOAuth},
   {"/integrations/connections/abc/test",c.RunTest},
   {"/integrations/connections/abc/tests",c.ListTests},
  } {
   recorder:=httptest.NewRecorder()
   request:=httptest.NewRequest(http.MethodPost,testCase.route,nil)
   testCase.handler(recorder,request)
   Expect(recorder.Code).To(Equal(http.StatusForbidden))
  }
 })
 It("denies connection creation without a verified principal",func(){
  c:=integrationcontroller.NewConnectionController(&serviceStub{})
  recorder:=httptest.NewRecorder()
  req:=httptest.NewRequest(http.MethodPost,"/integrations/connections",strings.NewReader(`{"provider":"GMAIL","ownerGid":"a"}`))
  req.Header.Set("X-Tenant-ID","spoofed")
  c.Create(recorder,req)
  Expect(recorder.Code).To(Equal(http.StatusForbidden))
 })
})

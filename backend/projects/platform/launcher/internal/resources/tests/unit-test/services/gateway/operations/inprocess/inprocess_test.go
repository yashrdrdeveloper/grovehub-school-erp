package gateway_test

import (
 "context"
 "errors"
 . "github.com/onsi/ginkgo/v2"
 . "github.com/onsi/gomega"
 "go.uber.org/mock/gomock"
 "grovehub/launcher/session"
 mock "grovehub/launcher/internal/resources/tests/mocks/gateway/inprocess"
 "grovehub/launcher/internal/src/services/gateway/operations/inprocess"
 "grovehub/launcher/internal/src/services/gateway/operations/inprocess/impl"
)

var authorizer *mock.MockAuthorizer

// @BeforeEach
func SetupAuthorizer(){
 authorizer=mock.NewMockAuthorizer(gomock.NewController(GinkgoT()))
}

// @Testing suite=LauncherInProcess
func InProcessSpecs(){
 Describe("Call",func(){
  It("invokes only after successful authorization",func(){
   principal:=session.Principal{AccountID:"account",SessionID:"session"}
   ctx:=session.WithPrincipal(context.Background(),principal)
   authorizer.EXPECT().Authorize(ctx,principal,"read").Return(nil)
   result,err:=inprocess.Invoke(ctx,impl.New(authorizer),"read",func(authorized context.Context)(string,error){
    Expect(authorized).To(Equal(ctx))
    return "result",nil
   })
   Expect(err).NotTo(HaveOccurred())
   Expect(result).To(Equal("result"))
  })
  It("does not invoke when authorization fails",func(){
   principal:=session.Principal{AccountID:"account",SessionID:"session"}
   ctx:=session.WithPrincipal(context.Background(),principal)
   denied:=errors.New("denied")
   authorizer.EXPECT().Authorize(ctx,principal,"write").Return(denied)
   called:=false
   err:=impl.New(authorizer).Call(ctx,"write",func(context.Context)error{called=true;return nil})
   Expect(err).To(MatchError(denied))
   Expect(called).To(BeFalse())
  })
 })
}

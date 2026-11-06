package providers_test

import (
 . "github.com/onsi/ginkgo/v2"
 . "github.com/onsi/gomega"

 "grovehub/services/integration/providers"
 mailfactory "grovehub/services/integration/providers/mailprovider/impl"
 "grovehub/services/integration/providers/mailprovider/providers/gmail"
 "grovehub/services/integration/providers/mailprovider/providers/outlook"
)

var _=Describe("Mail provider factory",func(){
 It("resolves providers and their declared capabilities without hard-coded service dispatch",func(){
  runtime,err:=mailfactory.NewFactory().Build(gmail.NewStrategy(nil),outlook.NewStrategy(nil))
  Expect(err).NotTo(HaveOccurred())
  Expect(runtime.Registry.List()).To(HaveLen(2))
  selection,err:=runtime.Resolver.Resolve("gmail",providers.CapabilityMailRead)
  Expect(err).NotTo(HaveOccurred())
  Expect(selection.Descriptor().ID).To(Equal("GMAIL"))
  _,err=runtime.Resolver.Resolve("OUTLOOK",providers.CapabilityMailWatch)
  Expect(err).To(MatchError(providers.ErrUnsupportedAction))
  Expect(runtime.Actions.Actions("GMAIL")).To(BeEmpty()) // no OAuth configuration
 })
 It("rejects duplicate provider strategies",func(){
  _,err:=mailfactory.NewFactory().Build(gmail.NewStrategy(nil),gmail.NewStrategy(nil))
  Expect(err).To(MatchError(providers.ErrDuplicateProvider))
 })
})

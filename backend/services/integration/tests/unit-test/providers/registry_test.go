package providers_test

import (
 . "github.com/onsi/ginkgo/v2"
 . "github.com/onsi/gomega"
 "grovehub/services/integration/providers"
 providerimpl "grovehub/services/integration/providers/impl"
 "grovehub/services/integration/providers/mailprovider/providers/gmail"
 "grovehub/services/integration/providers/mailprovider/providers/outlook"
)
var _ = Describe("Provider registry",func(){
 It("lists providers and resolves them case-insensitively",func(){
  registry,err:=providerimpl.NewRegistry(gmail.New(),outlook.New())
  Expect(err).NotTo(HaveOccurred())
  Expect(registry.List()).To(HaveLen(2))
  provider,err:=registry.Find("gmail")
  Expect(err).NotTo(HaveOccurred())
  Expect(provider.Descriptor().ID).To(Equal("GMAIL"))
 })
 It("rejects duplicate registrations",func(){
  _,err:=providerimpl.NewRegistry(gmail.New(),gmail.New())
  Expect(err).To(MatchError(providers.ErrDuplicateProvider))
 })
 It("rejects unsupported providers",func(){
  registry,err:=providerimpl.NewRegistry(gmail.New())
  Expect(err).NotTo(HaveOccurred())
  _,err=registry.Find("WHATSAPP")
  Expect(err).To(MatchError(providers.ErrUnknownProvider))
 })
})

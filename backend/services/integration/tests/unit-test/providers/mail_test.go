package providers_test

import (
 . "github.com/onsi/ginkgo/v2"
 . "github.com/onsi/gomega"
 "grovehub/services/integration/providers/gmail"
)
var _=Describe("Gmail mail normalization",func(){
 It("recognizes sent messages regardless of the other labels",func(){
  Expect(gmail.MailDirection([]string{"INBOX"})).To(Equal("INBOUND"))
  Expect(gmail.MailDirection([]string{"STARRED","SENT"})).To(Equal("OUTBOUND"))
  Expect(gmail.MailDirection([]string{"sent"})).To(Equal("OUTBOUND"))
 })
})

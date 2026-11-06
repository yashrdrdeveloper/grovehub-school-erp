package providers_test

import (
 "encoding/base64"
 "fmt"
 "strings"

 . "github.com/onsi/ginkgo/v2"
 . "github.com/onsi/gomega"
 "grovehub/services/integration/providers/mailprovider/providers/gmail"
)

var _=Describe("Gmail push payload",func(){
 It("decodes a Pub/Sub wrapped Gmail mailbox cursor",func(){
  encoded:=base64.StdEncoding.EncodeToString([]byte(`{"emailAddress":"ONE@GMAIL.COM","historyId":"12345"}`))
  payload:=fmt.Sprintf(`{"message":{"data":"%s","messageId":"id1"}}`,encoded)
  decoded,err:=gmail.DecodePush(strings.NewReader(payload))
  Expect(err).NotTo(HaveOccurred())
  Expect(decoded.EmailAddress).To(Equal("one@gmail.com"))
  Expect(decoded.HistoryID).To(Equal("12345"))
 })
 It("rejects invalid wrapped payloads",func(){
  _,err:=gmail.DecodePush(strings.NewReader(`{"message":{"data":"not base64"}}`))
  Expect(err).To(HaveOccurred())
 })
});

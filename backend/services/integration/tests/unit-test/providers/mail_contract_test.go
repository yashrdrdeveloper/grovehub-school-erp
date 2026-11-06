package providers_test

import (
 "bytes"
 "encoding/base64"
 "io"
 "mime"
 "mime/multipart"
 "net/mail"
 "strings"

 . "github.com/onsi/ginkgo/v2"
 . "github.com/onsi/gomega"
 mailcontract "grovehub/services/integration/providers/mailprovider/mailcontract"
)

var _=Describe("Provider-neutral mail MIME codec",func(){
 It("creates multipart mail with text HTML and attachments",func(){
  payload:=mailcontract.Send{
   To:[]string{"yashrdrdeveloper@gmail.com"},
   Cc:[]string{"copy@example.com"},
   Bcc:[]string{"hidden@example.com"},
   ReplyTo:[]string{"reply@example.com"},
   Subject:"GrooveHub Test",
   Body:mailcontract.Body{Text:"Hello plain",HTML:"<p>Hello HTML</p>"},
   Attachments:[]mailcontract.Attachment{{
    Filename:"test.txt",ContentType:"text/plain",
    ContentBase64:base64.StdEncoding.EncodeToString([]byte("hello attachment")),
   }},
  }
  raw,err:=mailcontract.BuildMIME("sender@gmail.com",payload)
  Expect(err).NotTo(HaveOccurred())
  msg,err:=mail.ReadMessage(bytes.NewReader(raw))
  Expect(err).NotTo(HaveOccurred())
  Expect(msg.Header.Get("From")).To(Equal("sender@gmail.com"))
  Expect(msg.Header.Get("To")).To(Equal("yashrdrdeveloper@gmail.com"))
  Expect(msg.Header.Get("Cc")).To(Equal("copy@example.com"))
  Expect(msg.Header.Get("Bcc")).To(Equal("hidden@example.com"))
  mediaType,params,err:=mime.ParseMediaType(msg.Header.Get("Content-Type"))
  Expect(err).NotTo(HaveOccurred())
  Expect(mediaType).To(Equal("multipart/mixed"))
  reader:=multipart.NewReader(msg.Body,params["boundary"])
  part,err:=reader.NextPart()
  Expect(err).NotTo(HaveOccurred())
  Expect(part.Header.Get("Content-Type")).To(ContainSubstring("multipart/alternative"))
  alternate,err:=io.ReadAll(part)
  Expect(err).NotTo(HaveOccurred())
  Expect(string(alternate)).To(ContainSubstring("text/plain"))
  Expect(string(alternate)).To(ContainSubstring("text/html"))
  file,err:=reader.NextPart()
  Expect(err).NotTo(HaveOccurred())
  Expect(file.FileName()).To(Equal("test.txt"))
  content,err:=io.ReadAll(file)
  Expect(err).NotTo(HaveOccurred())
  decoded,err:=base64.StdEncoding.DecodeString(strings.TrimSpace(string(content)))
  Expect(err).NotTo(HaveOccurred())
  Expect(string(decoded)).To(Equal("hello attachment"))
  _,err=reader.NextPart()
  Expect(err).To(Equal(io.EOF))
 })
 It("rejects spoofed header injection, invalid addresses and bad attachment encoding",func(){
  valid:=mailcontract.Send{
   To:[]string{"recipient@example.com"},Subject:"Test",
   Body:mailcontract.Body{Text:"message"},
  }
  _,err:=mailcontract.BuildMIME("sender@example.com",valid)
  Expect(err).NotTo(HaveOccurred())
  invalid:=valid
  invalid.To=[]string{"a@example.com\r\nBcc:bad@example.com"}
  Expect(invalid.Validate()).To(HaveOccurred())
  invalid=valid
  invalid.Subject="Hello\nBcc:bad@example.com"
  Expect(invalid.Validate()).To(HaveOccurred())
  invalid=valid
  invalid.Attachments=[]mailcontract.Attachment{{Filename:"a.txt",ContentType:"text/plain",ContentBase64:"not-base64"}}
  Expect(invalid.Validate()).To(HaveOccurred())
  _,err=mailcontract.BuildMIME("spoof@example.com\r\nBcc:bad@example.com",valid)
  Expect(err).To(HaveOccurred())
 })
 It("serializes text-only messages without attachments",func(){
  raw,err:=mailcontract.BuildMIME("sender@example.com",mailcontract.Send{
   To:[]string{"recipient@example.com"},Subject:"Test",Body:mailcontract.Body{Text:"Hello"},
  })
  Expect(err).NotTo(HaveOccurred())
  msg,err:=mail.ReadMessage(bytes.NewReader(raw))
  Expect(err).NotTo(HaveOccurred())
  Expect(msg.Header.Get("Content-Transfer-Encoding")).To(Equal("base64"))
  encoded,err:=io.ReadAll(msg.Body)
  Expect(err).NotTo(HaveOccurred())
  decoded,err:=base64.StdEncoding.DecodeString(strings.TrimSpace(string(encoded)))
  Expect(err).NotTo(HaveOccurred())
  Expect(string(decoded)).To(Equal("Hello"))
 })
})

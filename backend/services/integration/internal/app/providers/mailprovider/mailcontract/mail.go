package mailcontract

import (
 "bytes"
 "encoding/base64"
 "errors"
 "fmt"
 "io"
 "mime"
 "mime/multipart"
 "net/mail"
 "net/textproto"
 "strings"
)

const MaxAttachmentBytes = 3 << 20

// Send is provider-independent. The authenticated connection selects the sender.
// contentBase64 is standard RFC 4648 base64, not a provider-specific encoding.
type Send struct {
 To []string `json:"to"`
 Cc []string `json:"cc,omitempty"`
 Bcc []string `json:"bcc,omitempty"`
 ReplyTo []string `json:"replyTo,omitempty"`
 Subject string `json:"subject"`
 Body Body `json:"body"`
 Attachments []Attachment `json:"attachments,omitempty"`
}
type Body struct {
 Text string `json:"text,omitempty"`
 HTML string `json:"html,omitempty"`
}
type Attachment struct {
 Filename string `json:"filename"`
 ContentType string `json:"contentType"`
 ContentBase64 string `json:"contentBase64"`
 ContentID string `json:"contentId,omitempty"`
 Inline bool `json:"inline,omitempty"`
}
type AttachmentData struct {
 AttachmentID string `json:"attachmentId,omitempty"`
 Filename string `json:"filename,omitempty"`
 ContentType string `json:"contentType,omitempty"`
 Size int64 `json:"size"`
 ContentBase64 string `json:"contentBase64"`
}
type Message struct {
 ID string `json:"id"`
 ThreadID string `json:"threadId,omitempty"`
 Subject string `json:"subject,omitempty"`
 From string `json:"from,omitempty"`
 To []string `json:"to,omitempty"`
 Cc []string `json:"cc,omitempty"`
 Bcc []string `json:"bcc,omitempty"`
 ReplyTo []string `json:"replyTo,omitempty"`
 Body Body `json:"body"`
 Snippet string `json:"snippet,omitempty"`
 ReceivedAt string `json:"receivedAt,omitempty"`
 SentAt string `json:"sentAt,omitempty"`
 Labels []string `json:"labels,omitempty"`
 Attachments []AttachmentData `json:"attachments,omitempty"`
}
type Page struct {
 Messages []Message `json:"messages"`
 NextPageToken string `json:"nextPageToken,omitempty"`
}
var ErrInvalidMail = errors.New("invalid mail input")

func checkAddresses(values []string,required bool) error {
 if required&&len(values)==0{return ErrInvalidMail}
 if len(values)>100{return ErrInvalidMail}
 for _,address:=range values {
  if len(address)>320||strings.ContainsAny(address,"\r\n") {return ErrInvalidMail}
  parsed,err:=mail.ParseAddress(address)
  if err!=nil||parsed.Address!=address{return ErrInvalidMail}
 }
 return nil
}
func (s Send) Validate()error{
 for _,values:=range [][]string{s.To,s.Cc,s.Bcc,s.ReplyTo}{
  if err:=checkAddresses(values,false);err!=nil{return err}
 }
 if len(s.To)+len(s.Cc)+len(s.Bcc)==0{return ErrInvalidMail}
 if strings.TrimSpace(s.Subject)==""||len(s.Subject)>998||strings.ContainsAny(s.Subject,"\r\n"){return ErrInvalidMail}
 if s.Body.Text==""&&s.Body.HTML==""{return ErrInvalidMail}
 if len(s.Body.Text)+len(s.Body.HTML)>2<<20{return ErrInvalidMail}
 if len(s.Attachments)>20{return ErrInvalidMail}
 size:=0
 for _,a:=range s.Attachments {
  if a.Filename==""||len(a.Filename)>255||strings.ContainsAny(a.Filename,"\r\n")||
    strings.ContainsAny(a.ContentType,"\r\n")||a.ContentType==""||
    strings.ContainsAny(a.ContentID,"\r\n") {return ErrInvalidMail}
  if _,_,err:=mime.ParseMediaType(a.ContentType);err!=nil{return ErrInvalidMail}
  data,err:=base64.StdEncoding.DecodeString(a.ContentBase64)
  if err!=nil||len(data)==0{return ErrInvalidMail}
  size+=len(data)
  if size>MaxAttachmentBytes{return ErrInvalidMail}
 }
 return nil
}
func checkSender(sender string)error {
 parsed,err:=mail.ParseAddress(sender)
 if err!=nil||parsed.Address!=sender||strings.ContainsAny(sender,"\r\n"){return ErrInvalidMail}
 return nil
}
func encodeBodyPart(w *multipart.Writer, contentType, content string) error {
 header:=make(textproto.MIMEHeader)
 header.Set("Content-Type",contentType+"; charset=UTF-8")
 header.Set("Content-Transfer-Encoding","base64")
 part,err:=w.CreatePart(header)
 if err!=nil{return err}
 return writeEncoded(part,[]byte(content))
}
func writeEncoded(dst io.Writer,source []byte)error{
 // RFC 2045 base64 body lines have at most 76 characters, separated by CRLF.
 encoded:=base64.StdEncoding.EncodeToString(source)
 for len(encoded)>76 {
  if _,err:=io.WriteString(dst,encoded[:76]+"\r\n");err!=nil{return err}
  encoded=encoded[76:]
 }
 _,err:=io.WriteString(dst,encoded+"\r\n")
 return err
}
func mimeBody(s Send)(string,[]byte,error){
 if s.Body.Text!=""&&s.Body.HTML!=""{
  var output bytes.Buffer
  writer:=multipart.NewWriter(&output)
  if err:=encodeBodyPart(writer,"text/plain",s.Body.Text);err!=nil{return "",nil,err}
  if err:=encodeBodyPart(writer,"text/html",s.Body.HTML);err!=nil{return "",nil,err}
  if err:=writer.Close();err!=nil{return "",nil,err}
  return "multipart/alternative; boundary="+writer.Boundary(),output.Bytes(),nil
 }
 content:=s.Body.Text
 mediaType:="text/plain"
 if s.Body.HTML!=""{content=s.Body.HTML;mediaType="text/html"}
 return mediaType+"; charset=UTF-8",[]byte(content),nil
}
// BuildMIME serializes a structured mail request to RFC 5322/MIME bytes.
// Providers encode the resulting bytes according to their own API protocols.
func BuildMIME(sender string,s Send)([]byte,error){
 if err:=s.Validate();err!=nil{return nil,err}
 if err:=checkSender(sender);err!=nil{return nil,err}
 var message bytes.Buffer
 for _,h:=range []struct{name,value string}{
  {"From",sender},{"To",strings.Join(s.To,", ")},{"Cc",strings.Join(s.Cc,", ")},
  {"Bcc",strings.Join(s.Bcc,", ")},{"Reply-To",strings.Join(s.ReplyTo,", ")},
 } {
  if h.value!=""{fmt.Fprintf(&message,"%s: %s\r\n",h.name,h.value)}
 }
 fmt.Fprintf(&message,"Subject: %s\r\nMIME-Version: 1.0\r\n",mime.QEncoding.Encode("utf-8",s.Subject))
 contentType,body,err:=mimeBody(s)
 if err!=nil{return nil,err}
 if len(s.Attachments)==0{
  fmt.Fprintf(&message,"Content-Type: %s\r\n",contentType)
  if !strings.HasPrefix(contentType,"multipart/"){
   message.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
   if err:=writeEncoded(&message,body);err!=nil{return nil,err}
  }else{
   message.WriteString("\r\n")
   message.Write(body)
  }
  return message.Bytes(),nil
 }
 writer:=multipart.NewWriter(&message)
 fmt.Fprintf(&message,"Content-Type: multipart/mixed; boundary=%q\r\n\r\n",writer.Boundary())
 partHeaders:=make(textproto.MIMEHeader)
 partHeaders.Set("Content-Type",contentType)
 if !strings.HasPrefix(contentType,"multipart/"){partHeaders.Set("Content-Transfer-Encoding","base64")}
 part,err:=writer.CreatePart(partHeaders)
 if err!=nil{return nil,err}
 if strings.HasPrefix(contentType,"multipart/"){_,err=part.Write(body)}else{err=writeEncoded(part,body)}
 if err!=nil{return nil,err}
 for _,attachment:=range s.Attachments{
  data,_:=base64.StdEncoding.DecodeString(attachment.ContentBase64)
  h:=make(textproto.MIMEHeader)
  mediaType,params,_:=mime.ParseMediaType(attachment.ContentType)
  params["name"]=attachment.Filename
  h.Set("Content-Type",mime.FormatMediaType(mediaType,params))
  disposition:="attachment"
  if attachment.Inline{disposition="inline"}
  h.Set("Content-Disposition",mime.FormatMediaType(disposition,map[string]string{"filename":attachment.Filename}))
  h.Set("Content-Transfer-Encoding","base64")
  if attachment.ContentID!=""{h.Set("Content-ID","<"+strings.Trim(attachment.ContentID,"<>")+">")}
  part,err:=writer.CreatePart(h)
  if err!=nil{return nil,err}
  if err:=writeEncoded(part,data);err!=nil{return nil,err}
 }
 if err:=writer.Close();err!=nil{return nil,err}
 return message.Bytes(),nil
}

package dto

import (
    "encoding/base64"
    "errors"
    "mime"
    "net/mail"
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


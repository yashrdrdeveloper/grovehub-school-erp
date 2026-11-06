package impl

import (
 "encoding/base64"
 "encoding/json"
 "net/mail"
 "strings"

 "grovehub/services/integration/providers/mailprovider/mailcontract"
)

type gmailPart struct {
 MIMEType string `json:"mimeType"`
 Filename string `json:"filename"`
 Headers []struct {Name string `json:"name"`; Value string `json:"value"`} `json:"headers"`
 Body struct {Data string `json:"data"`; AttachmentID string `json:"attachmentId"`; Size int64 `json:"size"`} `json:"body"`
 Parts []gmailPart `json:"parts"`
}
func decodeGmailBody(raw string) string {
 if raw==""{return ""}
 decoded,err:=base64.RawURLEncoding.DecodeString(strings.TrimRight(raw,"="))
 if err!=nil||len(decoded)>2<<20{return ""}
 return string(decoded)
}
func applyGmailPart(dst *mailcontract.Message,p gmailPart){
 if len(p.Headers)>0{
  for _,h:=range p.Headers {
   switch strings.ToLower(h.Name){
   case "subject":dst.Subject=h.Value
   case "from":dst.From=h.Value
   case "to":dst.To=splitGmailAddresses(h.Value)
   case "cc":dst.Cc=splitGmailAddresses(h.Value)
   case "bcc":dst.Bcc=splitGmailAddresses(h.Value)
   case "reply-to":dst.ReplyTo=splitGmailAddresses(h.Value)
   }
  }
 }
 switch strings.ToLower(p.MIMEType){
 case "text/plain":if dst.Body.Text==""{dst.Body.Text=decodeGmailBody(p.Body.Data)}
 case "text/html":if dst.Body.HTML==""{dst.Body.HTML=decodeGmailBody(p.Body.Data)}
 }
 if p.Body.AttachmentID!=""{
  dst.Attachments=append(dst.Attachments,mailcontract.AttachmentData{
   AttachmentID:p.Body.AttachmentID,Filename:p.Filename,ContentType:p.MIMEType,Size:p.Body.Size,
  })
 }
 for _,child:=range p.Parts{applyGmailPart(dst,child)}
}
func splitGmailAddresses(raw string)[]string{
 // Preserve display-name forms and quoted commas using net/mail.
 list,err:=mail.ParseAddressList(raw)
 if err!=nil{return []string{raw}}
 result:=make([]string,0,len(list))
 for _,address:=range list{result=append(result,address.Address)}
 return result
}
func normalizeGmailMessage(m Message)mailcontract.Message{
 out:=mailcontract.Message{ID:m.ID,ThreadID:m.ThreadID,Snippet:m.Snippet,Labels:m.LabelIDs}
 if len(m.Payload)>0{
  var part gmailPart
  if json.Unmarshal(m.Payload,&part)==nil{applyGmailPart(&out,part)}
 }
 return out
}
func normalizeGmailPage(page MessagePage)mailcontract.Page{
 result:=mailcontract.Page{NextPageToken:page.NextPageToken,Messages:make([]mailcontract.Message,0,len(page.Messages))}
 for _,msg:=range page.Messages{result.Messages=append(result.Messages,normalizeGmailMessage(msg))}
 return result
}

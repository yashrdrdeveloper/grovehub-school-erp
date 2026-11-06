package impl

import (
 "encoding/json"
 "strings"

 "grovehub/services/integration/providers/mailprovider/mailcontract"
)
type graphAddress struct {
 EmailAddress struct {Address string `json:"address"`; Name string `json:"name"`} `json:"emailAddress"`
}
func graphAddressValue(raw json.RawMessage)string{
 var address graphAddress
 if json.Unmarshal(raw,&address)!=nil{return ""}
 return address.EmailAddress.Address
}
func graphAddressList(raw json.RawMessage)[]string{
 var addresses []graphAddress
 if json.Unmarshal(raw,&addresses)!=nil{return nil}
 values:=make([]string,0,len(addresses))
 for _,a:=range addresses{if a.EmailAddress.Address!=""{values=append(values,a.EmailAddress.Address)}}
 return values
}
func normalizeOutlookMessage(m Message)mailcontract.Message{
 result:=mailcontract.Message{
  ID:m.ID,ThreadID:m.ConversationID,Subject:m.Subject,From:graphAddressValue(m.From),
  To:graphAddressList(m.ToRecipients),Cc:graphAddressList(m.CcRecipients),
  Bcc:graphAddressList(m.BccRecipients),ReplyTo:graphAddressList(m.ReplyTo),
  ReceivedAt:m.ReceivedDateTime,SentAt:m.SentDateTime,
 }
 var body struct {ContentType string `json:"contentType"`; Content string `json:"content"`}
 if json.Unmarshal(m.Body,&body)==nil {
  if strings.EqualFold(body.ContentType,"html"){result.Body.HTML=body.Content}else{result.Body.Text=body.Content}
 }
 return result
}
func normalizeOutlookPage(page MessagePage)mailcontract.Page{
 result:=mailcontract.Page{NextPageToken:page.NextLink,Messages:make([]mailcontract.Message,0,len(page.Messages))}
 for _,message:=range page.Messages{result.Messages=append(result.Messages,normalizeOutlookMessage(message))}
 return result
}

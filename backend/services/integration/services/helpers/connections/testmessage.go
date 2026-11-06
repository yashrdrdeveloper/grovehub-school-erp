package connections

import (
 "encoding/base64"
 "errors"
 "fmt"
 "net/mail"

 "grovehub/services/integration/providers"
)

var ErrInvalidTestMailbox = errors.New("invalid test mailbox address")

// BuildSelfTestRaw produces a provider-neutral RFC 5322 MIME test message.
func BuildSelfTestRaw(mailbox,testID string)(string,error){
 address,err:=mail.ParseAddress(mailbox)
 if err!=nil||address.Address!=mailbox{
  return "",ErrInvalidTestMailbox
 }
 raw:=fmt.Sprintf(
  "From: <%s>\r\nTo: <%s>\r\nSubject: GrooveHub connection test\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\nGrooveHub outbound connectivity test. Test ID: %s\r\n",
  address.Address,address.Address,testID,
 )
 return base64.RawURLEncoding.EncodeToString([]byte(raw)),nil
}

// ProbeAction identifies provider-neutral actions; adapters implement behavior.
func ProbeAction(direction string)(providers.Action,error){
 switch direction{
 case "INBOUND":return providers.ActionMailList,nil
 case "OUTBOUND":return providers.ActionMailSend,nil
 default:return "",providers.ErrInvalidActionInput
 }
}

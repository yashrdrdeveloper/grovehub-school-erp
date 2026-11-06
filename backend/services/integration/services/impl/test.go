package impl

import (
 "context"
 "encoding/base64"
 "fmt"
 "net/mail"
 "strings"
 "time"

 "github.com/google/uuid"

 integrationresponse "grovehub/services/integration/dto/response"
 "grovehub/services/integration/entity"
)

func(s *ConnectionService)RunTest(ctx context.Context,tenantID,id,direction string)(integrationresponse.ConnectionTest,error){
 c,err:=s.repository.Get(ctx,tenantID,id)
 if err!=nil{return integrationresponse.ConnectionTest{},err}
 direction=strings.ToUpper(strings.TrimSpace(direction))
 if c.Status!=entity.StatusActive||c.ExternalID==nil{
  return integrationresponse.ConnectionTest{},ErrConnectionInactive
 }
 if c.Provider!="GMAIL"&&c.Provider!="OUTLOOK"{
  return integrationresponse.ConnectionTest{},ErrUnsupportedMailProvider
 }
 if direction!="INBOUND"&&direction!="OUTBOUND"{
  return integrationresponse.ConnectionTest{},ErrInvalidConnection
 }
 started:=time.Now().UTC()
 test:=entity.ConnectionTest{
  ID:uuid.NewString(),ConnectionID:id,Direction:direction,Status:"SUCCESS",StartedAt:started,
 }
 token,err:=s.validProviderToken(ctx,id,c.Provider)
 if err==nil{
  switch direction{
  case "INBOUND":
   switch c.Provider{
   case "GMAIL":
    if s.gmail==nil{return integrationresponse.ConnectionTest{},ErrOAuthUnavailable}
    err=s.gmail.InboundAccess(ctx,token)
   case "OUTLOOK":
    if s.outlook==nil{return integrationresponse.ConnectionTest{},ErrOAuthUnavailable}
    _,err=s.outlook.Messages(ctx,token.AccessToken,"")
   }
   if err==nil{test.Detail="Provider mailbox read capability verified; inbound delivery not verified"}
  case "OUTBOUND":
   var address *mail.Address
   address,err=mail.ParseAddress(*c.ExternalID)
   if err==nil&&address.Address!=*c.ExternalID{err=ErrInvalidConnection}
   if err==nil{
    message:=fmt.Sprintf("From: <%s>\r\nTo: <%s>\r\nSubject: GrooveHub connection test\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\nGrooveHub outbound connectivity test. Test ID: %s\r\n",address.Address,address.Address,test.ID)
    raw:=base64.RawURLEncoding.EncodeToString([]byte(message))
    switch c.Provider{
    case "GMAIL":
     if s.gmail==nil{return integrationresponse.ConnectionTest{},ErrOAuthUnavailable}
     err=s.gmail.Outbound(ctx,token,raw)
    case "OUTLOOK":
     if s.outlook==nil{return integrationresponse.ConnectionTest{},ErrOAuthUnavailable}
     _,err=s.outlook.SendRaw(ctx,token.AccessToken,raw)
    }
    if err==nil{test.Detail="Provider accepted self-addressed test message; delivery not confirmed"}
   }
  }
 }
 if err!=nil{
  test.Status="FAILED"
  test.Detail="Provider connection test failed; upstream error details redacted"
 }
 test.CompletedAt=time.Now().UTC()
 if err:=s.repository.SaveTest(ctx,test);err!=nil{return integrationresponse.ConnectionTest{},err}
 return toTest(test),nil
}

func(s *ConnectionService)ListTests(ctx context.Context,tenantID,id string)([]integrationresponse.ConnectionTest,error){
 if _,err:=s.repository.Get(ctx,tenantID,id);err!=nil{return nil,err}
 rows,err:=s.repository.ListTests(ctx,id)
 if err!=nil{return nil,err}
 result:=make([]integrationresponse.ConnectionTest,0,len(rows))
 for _,test:=range rows{result=append(result,toTest(test))}
 return result,nil
}
func toTest(test entity.ConnectionTest)integrationresponse.ConnectionTest{
 return integrationresponse.ConnectionTest{
  ID:test.ID,ConnectionID:test.ConnectionID,Direction:test.Direction,
  Status:test.Status,Detail:test.Detail,StartedAt:test.StartedAt,CompletedAt:test.CompletedAt,
 }
}

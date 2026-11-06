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
 if s.gmail==nil{return integrationresponse.ConnectionTest{},ErrOAuthUnavailable}
 c,err:=s.repository.Get(ctx,tenantID,id)
 if err!=nil{return integrationresponse.ConnectionTest{},err}
 direction=strings.ToUpper(strings.TrimSpace(direction))
 if c.Provider!="GMAIL"||c.Status!=entity.StatusActive||c.ExternalID==nil{return integrationresponse.ConnectionTest{},ErrConnectionInactive}
 if direction!="INBOUND"&&direction!="OUTBOUND"{return integrationresponse.ConnectionTest{},ErrInvalidConnection}
 started:=time.Now().UTC()
 test:=entity.ConnectionTest{
  ID:uuid.NewString(),ConnectionID:id,Direction:direction,Status:"SUCCESS",
  StartedAt:started,
 }
 token,err:=s.validToken(ctx,id)
 if err==nil{
  switch direction{
  case "INBOUND":
   err=s.gmail.InboundAccess(ctx,token)
   if err==nil{test.Detail="Gmail read permission verified; inbound delivery was not verified"}
  case "OUTBOUND":
   // Send a self-addressed test only; no arbitrary third-party recipients.
   var address *mail.Address
   address,err=mail.ParseAddress(*c.ExternalID)
   if err==nil && address.Address!=*c.ExternalID{err=ErrInvalidConnection}
   if err==nil{
    message:=fmt.Sprintf("From: <%s>\r\nTo: <%s>\r\nSubject: GrooveHub connection test\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\nGrooveHub outgoing test. Test ID: %s\r\n",address.Address,address.Address,test.ID)
    err=s.gmail.Outbound(ctx,token,base64.RawURLEncoding.EncodeToString([]byte(message)))
    if err==nil{test.Detail="Gmail API accepted self-addressed test message; delivery not confirmed"}
   }
  }
 }
 if err!=nil{
  test.Status="FAILED"
  // Avoid persisting upstream response details, OAuth credentials, or email content.
  test.Detail="Provider test failed; consult sanitized service telemetry"
 }
 test.CompletedAt=time.Now().UTC()
 if saveErr:=s.repository.SaveTest(ctx,test);saveErr!=nil{return integrationresponse.ConnectionTest{},saveErr}
 return toTest(test),nil
}
func(s *ConnectionService)ListTests(ctx context.Context,tenantID,id string)([]integrationresponse.ConnectionTest,error){
 if _,err:=s.repository.Get(ctx,tenantID,id);err!=nil{return nil,err}
 tests,err:=s.repository.ListTests(ctx,id)
 if err!=nil{return nil,err}
 results:=make([]integrationresponse.ConnectionTest,0,len(tests))
 for _,t:=range tests{results=append(results,toTest(t))}
 return results,nil
}
func toTest(t entity.ConnectionTest)integrationresponse.ConnectionTest{
 return integrationresponse.ConnectionTest{
 ID:t.ID,ConnectionID:t.ConnectionID,Direction:t.Direction,Status:t.Status,
 Detail:t.Detail,StartedAt:t.StartedAt,CompletedAt:t.CompletedAt,
 }
}

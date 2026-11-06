package impl

import (
 "context"
 "encoding/json"
 "errors"
 "strings"
 "time"

 "github.com/google/uuid"

 connectionconstants "grovehub/services/integration/constants/connections"
 integrationresponse "grovehub/services/integration/dto/response"
 "grovehub/services/integration/entity"
 "grovehub/services/integration/providers"
 connectionhelpers "grovehub/services/integration/services/helpers/connections"
)

// RunTest uses the registered provider strategy, without Gmail/Outlook branching.
func(s *ConnectionService)RunTest(ctx context.Context,tenantID,id,direction string)(integrationresponse.ConnectionTest,error){
 connection,err:=s.repository.Get(ctx,tenantID,id)
 if err!=nil{return integrationresponse.ConnectionTest{},err}

 direction=strings.ToUpper(strings.TrimSpace(direction))
 if connection.Status!=entity.StatusActive||connection.ExternalID==nil{
  return integrationresponse.ConnectionTest{},ErrConnectionInactive
 }
 action,err:=connectionhelpers.ProbeAction(direction)
 if err!=nil{return integrationresponse.ConnectionTest{},ErrInvalidConnection}

 started:=time.Now().UTC()
 test:=entity.ConnectionTest{
  ID:uuid.NewString(),
  ConnectionID:id,
  Direction:direction,
  Status:connectionconstants.TestSuccess,
  StartedAt:started,
 }
 payload:=json.RawMessage("{}")

 if direction==connectionconstants.DirectionOutbound{
  var encoded string
  encoded,err=connectionhelpers.BuildSelfTestRaw(*connection.ExternalID,test.ID)
  if err==nil{payload,err=json.Marshal(map[string]string{"raw":encoded})}
 }

 if err==nil{
  var token providers.OAuthToken
  token,err=s.validProviderToken(ctx,id,connection.Provider)
  if err==nil{
   if s.actions==nil{
    err=providers.ErrProviderUnavailable
   }else if token.AccessToken==""{
    err=errors.New("provider token missing")
   }else{
    _,err=s.actions.Execute(ctx,providers.Session{
     ConnectionID:id,
     ProviderID:connection.Provider,
     AccessToken:token.AccessToken,
    },action,payload)
   }
  }
 }
 if err!=nil{
  test.Status=connectionconstants.TestFailed
  test.Detail="Provider connection test failed; upstream details redacted"
 }else if direction==connectionconstants.DirectionInbound{
  test.Detail="Provider mailbox read capability verified; inbound delivery not verified"
 }else{
  test.Detail="Provider accepted self-addressed test message; delivery not confirmed"
 }

 test.CompletedAt=time.Now().UTC()
 if err:=s.repository.SaveTest(ctx,test);err!=nil{return integrationresponse.ConnectionTest{},err}
 return connectionhelpers.TestToResponse(test),nil
}

func(s *ConnectionService)ListTests(ctx context.Context,tenantID,id string)([]integrationresponse.ConnectionTest,error){
 if _,err:=s.repository.Get(ctx,tenantID,id);err!=nil{return nil,err}
 rows,err:=s.repository.ListTests(ctx,id)
 if err!=nil{return nil,err}
 result:=make([]integrationresponse.ConnectionTest,0,len(rows))
 for _,test:=range rows{result=append(result,connectionhelpers.TestToResponse(test))}
 return result,nil
}

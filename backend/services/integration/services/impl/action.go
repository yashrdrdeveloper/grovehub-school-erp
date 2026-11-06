package impl

import (
 "context"
 "encoding/json"
 "errors"

 "grovehub/services/integration/entity"
 "grovehub/services/integration/providers"
)

// ExecuteAction resolves the tenant-scoped provider binding before dispatching
// the optional action strategy. No message persistence occurs in Integration.
func(s *ConnectionService)ExecuteAction(ctx context.Context,tenantID,connectionID string,action providers.Action,input json.RawMessage)(any,error){
 c,err:=s.repository.Get(ctx,tenantID,connectionID)
 if err!=nil{return nil,err}
 if c.Status!=entity.StatusActive{return nil,ErrConnectionInactive}
 if s.actions==nil{return nil,providers.ErrProviderUnavailable}
 if len(s.actions.Actions(c.Provider))==0{return nil,providers.ErrProviderUnavailable}
 token,err:=s.validProviderToken(ctx,connectionID,c.Provider)
 if err!=nil{return nil,err}
 if token.AccessToken==""{return nil,errors.New("provider authentication unavailable")}
 return s.actions.Execute(ctx,providers.Session{
  ConnectionID:connectionID,ProviderID:c.Provider,AccessToken:token.AccessToken,
 },action,input)
}

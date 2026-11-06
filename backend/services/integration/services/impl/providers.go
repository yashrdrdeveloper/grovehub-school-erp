package impl

import (
 "context"

 integrationresponse "grovehub/services/integration/dto/response"
 "grovehub/services/integration/providers"
)

// Provider discovery does not create connections or execute provider actions.
func (s *ConnectionService) Providers(_ context.Context) []providers.Descriptor { return s.registry.List() }

func (s *ConnectionService) Provider(_ context.Context,id string)(integrationresponse.ProviderDetails,error){
 item,err:=s.registry.Find(id)
 if err!=nil{return integrationresponse.ProviderDetails{},err}
 descriptor:=item.Descriptor()
 _,oauthErr:=s.oauth.Find(descriptor.ID)
 return integrationresponse.ProviderDetails{
  Provider:descriptor,
  AvailableActions:s.actions.Actions(descriptor.ID),
  OAuthConfigured:oauthErr==nil,
 },nil
}

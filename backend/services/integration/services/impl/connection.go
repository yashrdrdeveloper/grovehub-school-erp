package impl

import (
 "context"
 "errors"
 "strings"
 "time"

 "github.com/google/uuid"

 coreid "grovehub/libs/core/id"
 integrationrequest "grovehub/services/integration/dto/request"
 integrationresponse "grovehub/services/integration/dto/response"
 "grovehub/services/integration/entity"
 "grovehub/services/integration/providers/gmail"
 "grovehub/services/integration/providers/outlook"
 infrastructureclient "grovehub/services/infrastructure/client"
 "grovehub/services/integration/providers"
 integrationrepository "grovehub/services/integration/repository"
 integrationservices "grovehub/services/integration/services"
)

var ErrInvalidConnection = errors.New("invalid integration connection")
var ErrConnectionNotReady = errors.New("provider authorization is not implemented")

type ConnectionService struct {
 repository integrationrepository.ConnectionRepository
 watchRepository integrationrepository.WatchRepository
 pushTopic string
 messaging infrastructureclient.MessagingClient
 registry providers.Registry
 actions providers.ActionRegistry
 oauth providers.OAuthRegistry
 gmail *gmail.Client
 outlook *outlook.Client
 cipher tokenCipher
}
func NewConnectionService(repository integrationrepository.ConnectionRepository,watchRepository integrationrepository.WatchRepository,registry providers.Registry,actions providers.ActionRegistry,oauth providers.OAuthRegistry, gmailClient *gmail.Client,outlookClient *outlook.Client, encryptionKey []byte,pushTopic string,messaging infrastructureclient.MessagingClient) (integrationservices.ConnectionService,error) {
 cipher,err:=newTokenCipher(encryptionKey)
 if err!=nil && (gmailClient!=nil||outlookClient!=nil){return nil,err}
 return &ConnectionService{repository:repository,watchRepository:watchRepository,pushTopic:pushTopic,messaging:messaging,registry:registry,actions:actions,oauth:oauth,gmail:gmailClient,outlook:outlookClient,cipher:cipher},nil
}
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
func (s *ConnectionService) Create(ctx context.Context,tenantID string,input integrationrequest.CreateConnection)(integrationresponse.Connection,error){
 if !isUUID(tenantID) || !isUUID(input.OwnerGID) {return integrationresponse.Connection{},ErrInvalidConnection}
 provider,err:=s.registry.Find(input.Provider)
 if err!=nil{return integrationresponse.Connection{},err}
 id,err:=coreid.NewUUIDv7()
 if err!=nil{return integrationresponse.Connection{},err}
 now:=time.Now().UTC()
 c:=entity.Connection{
  ID:id,TenantID:tenantID,OwnerGID:input.OwnerGID,
  Provider:provider.Descriptor().ID,Status:entity.StatusPendingSetup,
  CreatedAt:now,UpdatedAt:now,
 }
 if err=s.repository.Create(ctx,c);err!=nil{return integrationresponse.Connection{},err}
 return toConnection(c),nil
}
func (s *ConnectionService) Get(ctx context.Context,tenantID,id string)(integrationresponse.Connection,error){
 if !isUUID(tenantID)||!isUUID(id){return integrationresponse.Connection{},ErrInvalidConnection}
 c,err:=s.repository.Get(ctx,tenantID,id)
 if err!=nil{return integrationresponse.Connection{},err}
 return toConnection(c),nil
}
func (s *ConnectionService) List(ctx context.Context,tenantID,ownerGID string)([]integrationresponse.Connection,error){
 if !isUUID(tenantID){return nil,ErrInvalidConnection}
 rows,err:=s.repository.List(ctx,tenantID,ownerGID)
 if err!=nil{return nil,err}
 result:=make([]integrationresponse.Connection,0,len(rows))
 for _,c:=range rows{result=append(result,toConnection(c))}
 return result,nil
}
func (s *ConnectionService) Disconnect(ctx context.Context,tenantID,id string)(integrationresponse.Connection,error){
 if !isUUID(tenantID)||!isUUID(id){return integrationresponse.Connection{},ErrInvalidConnection}
 c,err:=s.repository.Disconnect(ctx,tenantID,id)
 if err!=nil{return integrationresponse.Connection{},err}
 return toConnection(c),nil
}
func toConnection(c entity.Connection) integrationresponse.Connection{
 return integrationresponse.Connection{
 ID:c.ID,TenantID:c.TenantID,OwnerGID:c.OwnerGID,Provider:c.Provider,
 Status:c.Status,ExternalID:c.ExternalID,LastTestedAt:c.LastTestedAt,LastTestStatus:c.LastTestStatus,CreatedAt:c.CreatedAt,UpdatedAt:c.UpdatedAt,
 }
}
func isUUID(raw string)bool{
 parsed,err:=uuid.Parse(strings.TrimSpace(raw))
 return err==nil && parsed.Version()==7
}

package services

import (
 "context"
 "encoding/json"

 integrationrequest "grovehub/services/integration/dto/request"
 integrationresponse "grovehub/services/integration/dto/response"
 "grovehub/services/integration/providers"
 "grovehub/services/integration/providers/gmail"
)

// ConnectionService owns provider connectivity and offers provider capabilities.
// It never persists email messages, threads, or consumer synchronization cursors.
type ConnectionService interface {
 Providers(ctx context.Context) []providers.Descriptor
 Provider(ctx context.Context,id string)(integrationresponse.ProviderDetails,error)
 Create(ctx context.Context,tenantID string,input integrationrequest.CreateConnection)(integrationresponse.Connection,error)
 Get(ctx context.Context,tenantID,id string)(integrationresponse.Connection,error)
 List(ctx context.Context,tenantID,ownerGID string)([]integrationresponse.Connection,error)
 Disconnect(ctx context.Context,tenantID,id string)(integrationresponse.Connection,error)
 BeginOAuth(ctx context.Context,tenantID,id string)(integrationresponse.OAuthStart,error)
 FinishOAuth(ctx context.Context,provider,code,state string)(integrationresponse.OAuthCompletion,error)
 RunTest(ctx context.Context,tenantID,id,direction string)(integrationresponse.ConnectionTest,error)
 ListTests(ctx context.Context,tenantID,id string)([]integrationresponse.ConnectionTest,error)
 Watch(ctx context.Context,tenantID,id string)(integrationresponse.Watch,error)
 ProcessPush(ctx context.Context,notification gmail.PushNotification)error
 RenewExpiring(ctx context.Context)error
 ExecuteAction(ctx context.Context,tenantID,connectionID string,action providers.Action,input json.RawMessage)(any,error)
}

package services

import (
 "context"
 integrationrequest "grovehub/services/integration/dto/request"
 integrationresponse "grovehub/services/integration/dto/response"
 "grovehub/services/integration/providers"
)

type ConnectionService interface {
 Providers(ctx context.Context) []providers.Descriptor
 Create(ctx context.Context,tenantID string,input integrationrequest.CreateConnection)(integrationresponse.Connection,error)
 Get(ctx context.Context,tenantID,id string)(integrationresponse.Connection,error)
 List(ctx context.Context,tenantID string)([]integrationresponse.Connection,error)
 Disconnect(ctx context.Context,tenantID,id string)(integrationresponse.Connection,error)
 BeginOAuth(ctx context.Context,tenantID,id string)(integrationresponse.OAuthStart,error)
 FinishOAuth(ctx context.Context,code,state string)(integrationresponse.OAuthCompletion,error)
 RunTest(ctx context.Context,tenantID,id,direction string)(integrationresponse.ConnectionTest,error)
 ListTests(ctx context.Context,tenantID,id string)([]integrationresponse.ConnectionTest,error)
 Sync(ctx context.Context,tenantID,id string)(integrationresponse.SyncResult,error)
}

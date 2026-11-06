package services

import (
 "context"

 integrationrequest "grovehub/services/integration/dto/request"
 integrationresponse "grovehub/services/integration/dto/response"
 "grovehub/services/integration/providers"
 "grovehub/services/integration/providers/gmail"
)

// ConnectionService owns provider connectivity and offers provider capabilities.
// It never persists email messages, threads, or consumer synchronization cursors.
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
 Watch(ctx context.Context,tenantID,id string)(integrationresponse.Watch,error)
 ProcessPush(ctx context.Context,notification gmail.PushNotification)error
 RenewExpiring(ctx context.Context)error
 MailProfile(ctx context.Context,tenantID,id string)(gmail.MailProfile,error)
 MailMessages(ctx context.Context,tenantID,id,page string)(gmail.MessagePage,error)
 MailMessage(ctx context.Context,tenantID,id,messageID string)(gmail.Message,error)
 MailHistory(ctx context.Context,tenantID,id,startHistoryID,page string)(gmail.HistoryPage,error)
 MailSend(ctx context.Context,tenantID,id,raw string)(gmail.SentMessage,error)
}

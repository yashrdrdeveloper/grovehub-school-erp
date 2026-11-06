package client

import (
 "context"
 "encoding/json"

 integrationrequest "grovehub/services/integration/dto/request"
 integrationresponse "grovehub/services/integration/dto/response"
 "grovehub/services/integration/providers"
)

// Client is the in-process service-to-service API, following Foundation's
// client/impl pattern. It does not bypass caller authorization decisions.
type Client interface{
 Providers(ctx context.Context)[]providers.Descriptor
 Provider(ctx context.Context,id string)(integrationresponse.ProviderDetails,error)
 CreateConnection(ctx context.Context,tenantID string,input integrationrequest.CreateConnection)(integrationresponse.Connection,error)
 GetConnection(ctx context.Context,tenantID,id string)(integrationresponse.Connection,error)
 ListConnections(ctx context.Context,tenantID,ownerGID string)([]integrationresponse.Connection,error)
 DisconnectConnection(ctx context.Context,tenantID,id string)(integrationresponse.Connection,error)
 ExecuteAction(ctx context.Context,tenantID,connectionID string,action providers.Action,input json.RawMessage)(any,error)
}

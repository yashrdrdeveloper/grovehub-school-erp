package impl

import (
 "context"
 "encoding/json"

 integrationclient "grovehub/services/integration/client"
 integrationrequest "grovehub/services/integration/dto/request"
 integrationresponse "grovehub/services/integration/dto/response"
 "grovehub/services/integration/providers"
 integrationservices "grovehub/services/integration/services"
)

type integrationClient struct{
 connections integrationservices.ConnectionService
 providers integrationservices.ProviderService
}

func NewClient(connections integrationservices.ConnectionService,providers integrationservices.ProviderService)integrationclient.Client{
 return &integrationClient{connections:connections,providers:providers}
}

func(c *integrationClient)Providers(ctx context.Context)[]providers.Descriptor{
 return c.providers.Providers(ctx)
}

func(c *integrationClient)Provider(ctx context.Context,id string)(integrationresponse.ProviderDetails,error){
 return c.providers.Provider(ctx,id)
}

func(c *integrationClient)CreateConnection(ctx context.Context,tenantID string,input integrationrequest.CreateConnection)(integrationresponse.Connection,error){
 return c.connections.Create(ctx,tenantID,input)
}

func(c *integrationClient)GetConnection(ctx context.Context,tenantID,id string)(integrationresponse.Connection,error){
 return c.connections.Get(ctx,tenantID,id)
}

func(c *integrationClient)ListConnections(ctx context.Context,tenantID,ownerGID string)([]integrationresponse.Connection,error){
 return c.connections.List(ctx,tenantID,ownerGID)
}

func(c *integrationClient)DisconnectConnection(ctx context.Context,tenantID,id string)(integrationresponse.Connection,error){
 return c.connections.Disconnect(ctx,tenantID,id)
}

func(c *integrationClient)ExecuteAction(ctx context.Context,tenantID,connectionID string,action providers.Action,input json.RawMessage)(any,error){
 return c.connections.ExecuteAction(ctx,tenantID,connectionID,action,input)
}

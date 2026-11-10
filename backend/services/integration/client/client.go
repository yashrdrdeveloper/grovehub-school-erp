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
// @Client integrationClient
type Client interface {
	// @Proxy providers.Providers
	Providers(ctx context.Context) []providers.Descriptor
	// @Proxy providers.Provider
	Provider(ctx context.Context, id string) (integrationresponse.ProviderDetails, error)
	// @Proxy connections.Create
	CreateConnection(ctx context.Context, tenantID string, input integrationrequest.CreateConnection) (integrationresponse.Connection, error)
	// @Proxy connections.Get
	GetConnection(ctx context.Context, tenantID, id string) (integrationresponse.Connection, error)
	// @Proxy connections.List
	ListConnections(ctx context.Context, tenantID, ownerGID string) ([]integrationresponse.Connection, error)
	// @Proxy connections.Disconnect
	DisconnectConnection(ctx context.Context, tenantID, id string) (integrationresponse.Connection, error)
	// @Proxy connections.ExecuteAction
	ExecuteAction(ctx context.Context, tenantID, connectionID string, action providers.Action, input json.RawMessage) (any, error)
}

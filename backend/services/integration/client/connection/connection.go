package connection

import (
    providerdto "grovehub/services/integration/internal/app/providers/dto"
    providerdto "grovehub/services/integration/internal/app/dto/provider"
	"context"
	"encoding/json"

	connectiondto "grovehub/services/integration/internal/app/dto/connection"
	
	"grovehub/services/integration/internal/app/providers"
)

// Client is the in-process service-to-service API, following Foundation's
// client/impl pattern. It does not bypass caller authorization decisions.
// @Client integrationClient
type Client interface {
	// @Proxy providers.Providers
	Providers(ctx context.Context) []providerdto.Descriptor
	// @Proxy providers.Provider
	Provider(ctx context.Context, id string) (providerdto.ProviderDetails, error)
	// @Proxy connections.Create
	CreateConnection(ctx context.Context, tenantID string, input connectiondto.CreateConnection) (connectiondto.Connection, error)
	// @Proxy connections.Get
	GetConnection(ctx context.Context, tenantID, id string) (connectiondto.Connection, error)
	// @Proxy connections.List
	ListConnections(ctx context.Context, tenantID, ownerGID string) ([]connectiondto.Connection, error)
	// @Proxy connections.Disconnect
	DisconnectConnection(ctx context.Context, tenantID, id string) (connectiondto.Connection, error)
	// @Proxy connections.ExecuteAction
	ExecuteAction(ctx context.Context, tenantID, connectionID string, action providerdto.Action, input json.RawMessage) (any, error)
}

// Client-facing aliases expose typed Integration contracts without exporting
// private DTO packages.
type (
    CreateConnection = connectiondto.CreateConnection
    Connection = connectiondto.Connection
    ConnectionTest = connectiondto.ConnectionTest
    OAuthStart = connectiondto.OAuthStart
    OAuthCompletion = connectiondto.OAuthCompletion
    Watch = connectiondto.Watch
    ProviderDetails = providerdto.ProviderDetails
    Action = providerdto.Action
    Descriptor = providerdto.Descriptor
)

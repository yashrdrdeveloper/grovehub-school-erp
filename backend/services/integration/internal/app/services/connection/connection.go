package connection

import (
	"context"
	"encoding/json"

	connectiondto "grovehub/services/integration/internal/app/dto/connection"
	
	"grovehub/services/integration/internal/app/providers"
	"grovehub/services/integration/internal/app/providers/mailprovider/providers/gmail"
)

// ConnectionService owns provider connectivity and offers provider capabilities.
// It never persists email messages, threads, or consumer synchronization cursors.
type ConnectionService interface {
	Create(ctx context.Context, tenantID string, input connectiondto.CreateConnection) (connectiondto.Connection, error)
	Get(ctx context.Context, tenantID, id string) (connectiondto.Connection, error)
	List(ctx context.Context, tenantID, ownerGID string) ([]connectiondto.Connection, error)
	Disconnect(ctx context.Context, tenantID, id string) (connectiondto.Connection, error)
	BeginOAuth(ctx context.Context, tenantID, id string) (connectiondto.OAuthStart, error)
	FinishOAuth(ctx context.Context, provider, code, state string) (connectiondto.OAuthCompletion, error)
	RunTest(ctx context.Context, tenantID, id, direction string) (connectiondto.ConnectionTest, error)
	ListTests(ctx context.Context, tenantID, id string) ([]connectiondto.ConnectionTest, error)
	Watch(ctx context.Context, tenantID, id string) (connectiondto.Watch, error)
	ProcessPush(ctx context.Context, notification gmail.PushNotification) error
	RenewExpiring(ctx context.Context) error
	ExecuteAction(ctx context.Context, tenantID, connectionID string, action providers.Action, input json.RawMessage) (any, error)
}

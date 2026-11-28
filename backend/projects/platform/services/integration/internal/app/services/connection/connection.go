package connection

import (
    providerdto "grovehub/services/integration/internal/app/providers/dto"
	"context"
	"encoding/json"

	connectiondto "grovehub/services/integration/internal/app/dto/connection"
	
	"grovehub/services/integration/internal/app/externalService/mail/gmail"
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
	ExecuteAction(ctx context.Context, tenantID, connectionID string, action providerdto.Action, input json.RawMessage) (any, error)
	ExecuteSystemMail(ctx context.Context, tenantID string, input json.RawMessage) (any, error)
	ConfigureSystemMailSender(ctx context.Context, tenantID string, connectionID string) error
}

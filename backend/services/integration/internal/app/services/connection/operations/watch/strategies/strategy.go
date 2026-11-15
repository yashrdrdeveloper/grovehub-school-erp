package strategies

import (
    "context"
    connectiondto "grovehub/services/integration/internal/app/dto/connection"
    "grovehub/services/integration/internal/app/providers/mailprovider/providers/gmail"
)

// @Strategy
type Strategy interface {
    Watch(ctx context.Context, tenantID, id string) (connectiondto.Watch,error)
    ProcessPush(ctx context.Context, notification gmail.PushNotification) error
    RenewExpiring(ctx context.Context) error
}

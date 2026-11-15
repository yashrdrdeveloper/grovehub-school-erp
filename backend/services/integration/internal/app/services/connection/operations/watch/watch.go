package watch

import (
    "context"
    connectiondto "grovehub/services/integration/internal/app/dto/connection"
    "grovehub/services/integration/internal/app/providers/mailprovider/providers/gmail"
    "grovehub/services/integration/internal/app/services/connection/operations/watch/strategies"
)

type WatchOperation struct{ strategy strategies.Strategy }
func New(strategy strategies.Strategy) *WatchOperation {return &WatchOperation{strategy:strategy}}
func (operation *WatchOperation) Watch(ctx context.Context, tenantID, id string) (connectiondto.Watch,error) {return operation.strategy.Watch(ctx,tenantID,id)}
func (operation *WatchOperation) ProcessPush(ctx context.Context, notification gmail.PushNotification) error {return operation.strategy.ProcessPush(ctx,notification)}
func (operation *WatchOperation) RenewExpiring(ctx context.Context) error {return operation.strategy.RenewExpiring(ctx)}

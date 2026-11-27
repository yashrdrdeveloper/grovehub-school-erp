package disconnect

import (
    "context"
    connectiondto "grovehub/services/integration/internal/app/dto/connection"
    "grovehub/services/integration/internal/app/services/connection/operations/disconnect/strategies"
)

type DisconnectOperation struct{strategy strategies.Strategy}

func New(strategy strategies.Strategy) *DisconnectOperation {return &DisconnectOperation{strategy:strategy}}

func (operation *DisconnectOperation) Disconnect(ctx context.Context, tenantID, id string) (connectiondto.Connection,error) {return operation.strategy.Execute(ctx,tenantID,id)}

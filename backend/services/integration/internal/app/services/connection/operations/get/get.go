package get

import (
    "context"
    connectiondto "grovehub/services/integration/internal/app/dto/connection"
    "grovehub/services/integration/internal/app/services/connection/operations/get/strategies"
)

type GetOperation struct{strategy strategies.Strategy}
func New(strategy strategies.Strategy) *GetOperation {return &GetOperation{strategy:strategy}}
func (operation *GetOperation) Get(ctx context.Context, tenantID, id string) (connectiondto.Connection,error) {return operation.strategy.Execute(ctx,tenantID,id)}
